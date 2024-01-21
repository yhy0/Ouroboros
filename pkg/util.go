package pkg

import (
    "math/rand"
    "time"
)

/**
   @author yhy
   @since 2024/1/20
   @desc //TODO
**/

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()?"

func RandomString() string {
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, 10)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
