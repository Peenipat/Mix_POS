
package aws

import (
    "context"
    "fmt"
    "mime/multipart"
    "os"
    "path/filepath"
    "time"
    "crypto/rand"
    "encoding/hex"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func generateRandomFilename(original string) string {
    ext := filepath.Ext(original)
    b := make([]byte, 8)
    _, _ = rand.Read(b)
    randomStr := hex.EncodeToString(b)
    return fmt.Sprintf("%d_%s%s", time.Now().Unix(), randomStr, ext)
}

func UploadToS3(fileHeader *multipart.FileHeader, keyPrefix string) (string, string, error) {
    f, err := fileHeader.Open()
    if err != nil {
        return "", "", err
    }
    defer f.Close()

    filename := generateRandomFilename(fileHeader.Filename)
    key := filepath.Join(keyPrefix, filename) 

    // อัปโหลดขึ้น S3
    _, err = S3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
        Key:    aws.String(key),
        Body:   f,
        ACL:    types.ObjectCannedACLPublicRead,
    })
    if err != nil {
        return "", "", err
    }

    return keyPrefix, filename, nil
}

