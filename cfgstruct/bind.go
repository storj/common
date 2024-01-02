// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cfgstruct

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"storj.io/common/version"
)

const (
	// AnySource is a source annotation for config values that can come from
	// a flag or file.
	AnySource = "any"

	// FlagSource is a source annotation for config values that just come from
	// flags (i.e. are never persisted to file).
	FlagSource = "flag"

	// BasicHelpAnnotationName is the name of the annotation used to indicate
	// a flag should be included in basic usage/help.
	BasicHelpAnnotationName = "basic-help"
)

var (
	allSources = []string{
		AnySource,
		FlagSource,
	}
)

// BindOpt is an option for the Bind method.
type BindOpt struct {
	isDev   *bool
	isTest  *bool
	isSetup *bool
	varfn   func(vars map[string]confVar)
	prefix  string
}

// ConfDir sets variables for default options called $CONFDIR.
func ConfDir(path string) BindOpt {
	return ConfigVar("CONFDIR", filepath.Clean(os.ExpandEnv(path)))
}

// IdentityDir sets a variable for the default option called $IDENTITYDIR.
func IdentityDir(path string) BindOpt {
	return ConfigVar("IDENTITYDIR", filepath.Clean(os.ExpandEnv(path)))
}

// ConfigVar sets a variable for the default option called name.
func ConfigVar(name, val string) BindOpt {
	name = strings.ToUpper(name)
	return BindOpt{varfn: func(vars map[string]confVar) {
		vars[name] = confVar{val: val, nested: false}
	}}
}

// SetupMode issues the bind in a mode where it does not ignore fields with the
// `setup:"true"` tag.
func SetupMode() BindOpt {
	setup := true
	return BindOpt{isSetup: &setup}
}

// UseDevDefaults forces the bind call to use development defaults unless
// something else is provided as a subsequent option.
// Without a specific defaults setting, Bind will default to determining which
// defaults to use based on version.Build.Release.
func UseDevDefaults() BindOpt {
	dev := true
	test := false
	return BindOpt{isDev: &dev, isTest: &test}
}

// UseReleaseDefaults forces the bind call to use release defaults unless
// something else is provided as a subsequent option.
// Without a specific defaults setting, Bind will default to determining which
// defaults to use based on version.Build.Release.
func UseReleaseDefaults() BindOpt {
	dev := false
	test := false
	return BindOpt{isDev: &dev, isTest: &test}
}

// UseTestDefaults forces the bind call to use test defaults unless
// something else is provided as a subsequent option.
// Without a specific defaults setting, Bind will default to determining which
// defaults to use based on version.Build.Release.
func UseTestDefaults() BindOpt {
	dev := false
	test := true
	return BindOpt{isDev: &dev, isTest: &test}
}

type confVar struct {
	val    string
	nested bool
}

// Bind sets flags on a FlagSet that match the configuration struct
// 'config'. This works by traversing the config struct using the 'reflect'
// package.
func Bind(flags FlagSet, config interface{}, opts ...BindOpt) {
	bind(flags, config, opts...)
}

func bind(flags FlagSet, config interface{}, opts ...BindOpt) {
	ptrtype := reflect.TypeOf(config)
	if ptrtype.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("invalid config type: %#v. Expecting pointer to struct.", config))
	}
	isDev := !version.Build.Release
	isTest := false
	setupCommand := false
	vars := map[string]confVar{}
	prefix := ""
	for _, opt := range opts {
		if opt.varfn != nil {
			opt.varfn(vars)
		}
		if opt.isDev != nil {
			isDev = *opt.isDev
		}
		if opt.isTest != nil {
			isTest = *opt.isTest
		}
		if opt.isSetup != nil {
			setupCommand = *opt.isSetup
		}
		if opt.prefix != "" {
			prefix = strings.TrimSuffix(opt.prefix, ".") + "."
		}
	}

	bindConfig(flags, prefix, reflect.ValueOf(config).Elem(), vars, setupCommand, false, isDev, isTest)
}

func bindConfig(flags FlagSet, prefix string, val reflect.Value, vars map[string]confVar, setupCommand, setupStruct bool, isDev, isTest bool) {
	if val.Kind() != reflect.Struct {
		panic(fmt.Sprintf("invalid config type: %#v. Expecting struct.", val.Interface()))
	}
	typ := val.Type()
	resolvedVars := make(map[string]string, len(vars))
	{
		structpath := strings.ReplaceAll(prefix, ".", string(filepath.Separator))
		for k, v := range vars {
			if !v.nested {
				resolvedVars[k] = v.val
				continue
			}
			resolvedVars[k] = filepath.Join(v.val, structpath)
		}
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldval := val.Field(i)
		flagname := hyphenate(snakeCase(field.Name))

		if field.Tag.Get("noprefix") != "true" {
			flagname = prefix + flagname
		}

		onlyForSetup := (field.Tag.Get("setup") == "true") || setupStruct
		// ignore setup params for non setup commands
		if !setupCommand && onlyForSetup {
			continue
		}

		if !fieldval.CanAddr() {
			panic(fmt.Sprintf("cannot addr field %s in %s", field.Name, typ))
		}

		fieldref := fieldval.Addr()
		if !fieldref.CanInterface() {
			panic(fmt.Sprintf("cannot get interface of field %s in %s", field.Name, typ))
		}

		fieldaddr := fieldref.Interface()
		if fieldvalue, ok := fieldaddr.(pflag.Value); ok {
			help := field.Tag.Get("help")
			def := getDefault(field.Tag, isTest, isDev, flagname)

			if field.Tag.Get("internal") == "true" {
				if def != "" {
					panic(fmt.Sprintf("unapplicable default value set for internal flag: %s", flagname))
				}
				continue
			}

			err := fieldvalue.Set(def)
			if err != nil {
				panic(fmt.Sprintf("invalid default value for %s: %#v, %v", flagname, def, err))
			}
			flags.Var(fieldvalue, flagname, help)

			markHidden := false
			if onlyForSetup {
				SetBoolAnnotation(flags, flagname, "setup", true)
			}
			if field.Tag.Get("user") == "true" {
				SetBoolAnnotation(flags, flagname, "user", true)
			}
			if field.Tag.Get("hidden") == "true" {
				markHidden = true
				SetBoolAnnotation(flags, flagname, "hidden", true)
			}
			if field.Tag.Get("deprecated") == "true" {
				markHidden = true
				SetBoolAnnotation(flags, flagname, "deprecated", true)
			}
			if source := field.Tag.Get("source"); source != "" {
				setSourceAnnotation(flags, flagname, source)
			}
			if markHidden {
				err := flags.MarkHidden(flagname)
				if err != nil {
					panic(fmt.Sprintf("mark hidden failed %s: %v", flagname, err))
				}
			}
			continue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			if field.Anonymous {
				bindConfig(flags, prefix, fieldval, vars, setupCommand, onlyForSetup, isDev, isTest)
			} else {
				bindConfig(flags, flagname+".", fieldval, vars, setupCommand, onlyForSetup, isDev, isTest)
			}
		case reflect.Array:
			digits := len(fmt.Sprint(fieldval.Len()))
			for j := 0; j < fieldval.Len(); j++ {
				padding := strings.Repeat("0", digits-len(fmt.Sprint(j)))
				bindConfig(flags, fmt.Sprintf("%s.%s%d.", flagname, padding, j), fieldval.Index(j), vars, setupCommand, onlyForSetup, isDev, isTest)
			}
		default:
			help := field.Tag.Get("help")
			def := getDefault(field.Tag, isTest, isDev, flagname)

			if field.Tag.Get("internal") == "true" {
				if def != "" {
					panic(fmt.Sprintf("unapplicable default value set for internal flag: %s", flagname))
				}
				continue
			}

			def = expand(resolvedVars, def)

			fieldaddr := fieldval.Addr().Interface()
			check := func(err error) {
				if err != nil {
					panic(fmt.Sprintf("invalid default value for %s: %#v", flagname, def))
				}
			}
			switch field.Type {
			case reflect.TypeOf(int(0)):
				val, err := strconv.ParseInt(def, 0, strconv.IntSize)
				check(err)
				flags.IntVar(fieldaddr.(*int), flagname, int(val), help)
			case reflect.TypeOf(int64(0)):
				val, err := strconv.ParseInt(def, 0, 64)
				check(err)
				flags.Int64Var(fieldaddr.(*int64), flagname, val, help)
			case reflect.TypeOf(uint(0)):
				val, err := strconv.ParseUint(def, 0, strconv.IntSize)
				check(err)
				flags.UintVar(fieldaddr.(*uint), flagname, uint(val), help)
			case reflect.TypeOf(uint64(0)):
				val, err := strconv.ParseUint(def, 0, 64)
				check(err)
				flags.Uint64Var(fieldaddr.(*uint64), flagname, val, help)
			case reflect.TypeOf(time.Duration(0)):
				val, err := time.ParseDuration(def)
				check(err)
				flags.DurationVar(fieldaddr.(*time.Duration), flagname, val, help)
			case reflect.TypeOf(float64(0)):
				val, err := strconv.ParseFloat(def, 64)
				check(err)
				flags.Float64Var(fieldaddr.(*float64), flagname, val, help)
			case reflect.TypeOf(string("")):
				if field.Tag.Get("path") == "true" {
					// NB: conventionally unix path separators are used in default values
					def = filepath.FromSlash(def)
				}
				flags.StringVar(fieldaddr.(*string), flagname, def, help)
			case reflect.TypeOf(bool(false)):
				val, err := strconv.ParseBool(def)
				check(err)
				flags.BoolVar(fieldaddr.(*bool), flagname, val, help)
			case reflect.TypeOf([]string(nil)):
				// allow either a single string, or comma separated values for defaults
				defaultValues := []string{}
				if def != "" {
					defaultValues = strings.Split(def, ",")
				}
				flags.StringSliceVar(fieldaddr.(*[]string), flagname, defaultValues, help)
			default:
				panic(fmt.Sprintf("invalid field type: %s", field.Type.String()))
			}
			if onlyForSetup {
				SetBoolAnnotation(flags, flagname, "setup", true)
			}
			if field.Tag.Get("user") == "true" {
				SetBoolAnnotation(flags, flagname, "user", true)
			}
			if field.Tag.Get(BasicHelpAnnotationName) == "true" {
				SetBoolAnnotation(flags, flagname, BasicHelpAnnotationName, true)
			}

			markHidden := false
			if field.Tag.Get("hidden") == "true" {
				markHidden = true
				SetBoolAnnotation(flags, flagname, "hidden", true)
			}
			if field.Tag.Get("deprecated") == "true" {
				markHidden = true
				SetBoolAnnotation(flags, flagname, "deprecated", true)
			}
			if source := field.Tag.Get("source"); source != "" {
				setSourceAnnotation(flags, flagname, source)
			}
			if markHidden {
				err := flags.MarkHidden(flagname)
				if err != nil {
					panic(fmt.Sprintf("mark hidden failed %s: %v", flagname, err))
				}
			}
		}
	}
}

func getDefault(tag reflect.StructTag, isTest, isDev bool, flagname string) string {
	var order []string
	var opposites []string
	if isTest {
		order = []string{"testDefault", "devDefault", "default"}
		opposites = []string{"releaseDefault"}
	} else if isDev {
		order = []string{"devDefault", "default"}
		opposites = []string{"releaseDefault", "testDefault"}
	} else {
		order = []string{"releaseDefault", "default"}
		opposites = []string{"devDefault", "testDefault"}
	}

	for _, name := range order {
		if val, ok := tag.Lookup(name); ok {
			return val
		}
	}

	for _, name := range opposites {
		if _, ok := tag.Lookup(name); ok {
			panic(fmt.Sprintf("%q missing but %q defined for %v", order[0], name, flagname))
		}
	}

	return ""
}

func setSourceAnnotation(flagset interface{}, name, source string) {
	switch source {
	case AnySource:
	case FlagSource:
	default:
		panic(fmt.Sprintf("invalid source annotation %q for %s: must be one of %q", source, name, allSources))
	}

	setStringAnnotation(flagset, name, "source", source)
}

func setStringAnnotation(flagset interface{}, name, key, value string) {
	flags, ok := flagset.(*pflag.FlagSet)
	if !ok {
		return
	}

	err := flags.SetAnnotation(name, key, []string{value})
	if err != nil {
		panic(fmt.Sprintf("unable to set %s annotation for %s: %v", key, name, err))
	}
}

// SetBoolAnnotation sets an annotation (if it can) on flagset with a value of []string{"true|false"}.
func SetBoolAnnotation(flagset interface{}, name, key string, value bool) {
	flags, ok := flagset.(*pflag.FlagSet)
	if !ok {
		return
	}

	err := flags.SetAnnotation(name, key, []string{strconv.FormatBool(value)})
	if err != nil {
		panic(fmt.Sprintf("unable to set %s annotation for %s: %v", key, name, err))
	}
}

func expand(vars map[string]string, val string) string {
	return os.Expand(val, func(key string) string { return vars[key] })
}

// FindConfigDirParam returns '--config-dir' param from os.Args (if exists).
func FindConfigDirParam() string {
	return FindFlagEarly("config-dir")
}

// FindIdentityDirParam returns '--identity-dir' param from os.Args (if exists).
func FindIdentityDirParam() string {
	return FindFlagEarly("identity-dir")
}

// FindDefaultsParam returns '--defaults' param from os.Args (if it exists).
func FindDefaultsParam() string {
	return FindFlagEarly("defaults")
}

// FindFlagEarly retrieves the value of a flag before `flag.Parse` has been called.
func FindFlagEarly(flagName string) string {
	// workaround to have early access to 'dir' param
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, fmt.Sprintf("--%s=", flagName)) {
			return strings.TrimPrefix(arg, fmt.Sprintf("--%s=", flagName))
		} else if arg == fmt.Sprintf("--%s", flagName) && i < len(os.Args)-1 {
			return os.Args[i+1]
		}
	}
	return ""
}

// SetupFlag sets up flags that are needed before `flag.Parse` has been called.
func SetupFlag(log *zap.Logger, cmd *cobra.Command, dest *string, name, value, usage string) {
	if foundValue := FindFlagEarly(name); foundValue != "" {
		value = foundValue
	}
	cmd.PersistentFlags().StringVar(dest, name, value, usage)
	if cmd.PersistentFlags().SetAnnotation(name, "setup", []string{"true"}) != nil {
		log.Error("Failed to set 'setup' annotation", zap.String("Flag", name))
	}
}

// DefaultsType returns the type of defaults (release/dev) this binary should use.
func DefaultsType() string {
	// define a flag so that the flag parsing system will be happy.
	defaults := strings.ToLower(FindDefaultsParam())
	if defaults != "" {
		return defaults
	}
	if version.Build.Release {
		return "release"
	}
	return "dev"
}

// Prefix defines the used prefix, where configs are bound to.
func Prefix(prefix string) BindOpt {
	return BindOpt{
		prefix: prefix,
	}
}

// DefaultsFlag sets up the defaults=dev/release flag options, which is needed
// before `flag.Parse` has been called.
func DefaultsFlag(cmd *cobra.Command) BindOpt {
	// define a flag so that the flag parsing system will be happy.
	defaults := DefaultsType()

	// we're actually going to ignore this flag entirely and parse the commandline
	// arguments early instead
	_ = cmd.PersistentFlags().String("defaults", defaults,
		"determines which set of configuration defaults to use. can either be 'dev' or 'release'")
	setSourceAnnotation(cmd.PersistentFlags(), "defaults", FlagSource)

	switch defaults {
	case "dev":
		return UseDevDefaults()
	case "release":
		return UseReleaseDefaults()
	case "test":
		return UseTestDefaults()
	default:
		panic(fmt.Sprintf("unsupported defaults value %q", FindDefaultsParam()))
	}
}
