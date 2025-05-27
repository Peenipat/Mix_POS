
package aws

import (
    "context"
    "fmt"
    "mime/multipart"
    "os"
    "path/filepath"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

// uploadToS3 รับ FileHeader แล้วคืน public URL กับชื่อไฟล์
func UploadToS3(fileHeader *multipart.FileHeader) (string, string, error) {
    f, err := fileHeader.Open()
    if err != nil {
        return "", "", err
    }
    defer f.Close()

    // ตั้งชื่อไฟล์ไม่ซ้ำ
    base := filepath.Base(fileHeader.Filename)
    filename := fmt.Sprintf("%d_%s", time.Now().Unix(), base)
    key := "avatars/" + filename

    out, err := S3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
        Key:    aws.String(key),
        Body:   f,
    })
    if err != nil {
        return "", "", err
    }

    return out.Location, filename, nil
}
