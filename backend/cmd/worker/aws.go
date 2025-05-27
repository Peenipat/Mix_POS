
package aws

import (
    "context"
    "log"
    "os"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

var S3Uploader *manager.Uploader

func InitAWS() {
    // โหลด config จาก ENV: AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(os.Getenv("AWS_REGION")),
    )
    if err != nil {
        log.Fatalf("unable to load AWS SDK config: %v", err)
    }

    client := s3.NewFromConfig(cfg)
    S3Uploader = manager.NewUploader(client)
}
