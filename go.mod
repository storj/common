module storj.io/common

go 1.13

// force specific versions for minio
require github.com/btcsuite/btcutil v0.0.0-20180706230648-ab6388e0c60a

replace google.golang.org/grpc => github.com/storj/grpc-go v1.23.1-0.20190918084400-1c4561bf5127

require (
	github.com/calebcase/tmpfile v1.0.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/minio/sha256-simd v0.0.0-20190328051042-05b4dd3047e5
	github.com/pkg/errors v0.8.1 // indirect
	github.com/skyrings/skyring-common v0.0.0-20160929130248-d1c0bb1cbd5e
	github.com/spacemonkeygo/monkit/v3 v3.0.1
	github.com/stretchr/testify v1.3.0
	github.com/zeebo/errs v1.2.2
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20200107144601-ef85f5a75ddf
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20190716160619-c506a9f90610 // indirect
	google.golang.org/grpc v1.23.1
	storj.io/drpc v0.0.7-0.20191115031725-2171c57838d2
)
