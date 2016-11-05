package dirsync

import (
    "fmt"
    "os"
    "encoding/hex"
    "io"
    "crypto/md5"
    "strings"
    "path/filepath"
    "github.com/Varjelus/kopsa"
)

const mega = 1000000

// Errors
const ErrNotADirectory = "Not a directory"

func identicalFiles(path1, path2 string) (bool, error) {
    buffer := make([]byte, 10 * mega)
    h1 := md5.New()
    h2 := md5.New()

    f1, err := os.Open(path1)
    if err != nil {
        return false, err
    }
    defer f1.Close()

    f2, err := os.Open(path2)
    if err != nil {
        return false, err
    }
    defer f2.Close()

    if _, err := io.CopyBuffer(h1, f1, buffer); err != nil {
        return false, err
    }

    if _, err := io.CopyBuffer(h2, f2, buffer); err != nil {
        return false, err
    }

    if hex.EncodeToString(h1.Sum(nil)) == hex.EncodeToString(h2.Sum(nil)) {
        return true, nil
    } else {
        return false, nil
    }

    return false, nil
}

func Sync(src, dest string) (err error) {
    // Make paths absolute
    src, err = filepath.Abs(src)
    if err != nil {
        return
    }
    dest, err = filepath.Abs(dest)
    if err != nil {
        return
    }

    // Create target directory if needed
    info, err := os.Stat(src)
    if err != nil {
        return
    }
    if !info.IsDir() {
        return fmt.Errorf(ErrNotADirectory)
    }
    if err := os.MkdirAll(dest, info.Mode()); err != nil {
        return err
    }

    // Delete extra files
    err = filepath.Walk(dest, func(path string, dfi os.FileInfo, err error) error {
        relPath := strings.TrimPrefix(path, dest)
        srcPath := filepath.Join(src, relPath)

        if dfi.IsDir() {
            return nil
        }

        if sfi, err := os.Stat(srcPath); err != nil {
            if os.IsNotExist(err) {
                // Delete
                if err := os.Remove(path); err != nil {
                    return err
                }
            } else {
                return err
            }
        } else {
            // Files with same names
            // First compare sizes
            if dfi.Size() == sfi.Size() {
                // If sizes match, compare checksums
                identical, err := identicalFiles(path, srcPath)
                if err != nil {
                    return err
                }
                if !identical {
                    if err := os.Remove(path); err != nil {
                        return err
                    }
                }
            } else {
                if err := os.Remove(path); err != nil {
                    return err
                }
            }
        }

        return nil
    })
    if err != nil {
        return
    }

    // Copy new files over
    err = filepath.Walk(src, func(path string, sfi os.FileInfo, err error) error {
        relPath := strings.TrimPrefix(path, src)
        destPath := filepath.Join(dest, relPath)

        if sfi.IsDir() {
            info, err := os.Stat(path)
            if err != nil {
                return err
            }
            if err := os.MkdirAll(destPath, info.Mode()); err != nil {
                return err
            }
            return nil
        }

        if _, err := os.Stat(destPath); err != nil {
            if !os.IsNotExist(err) {
                return err
            }
            _, err := kopsa.Copy(destPath, path)
            if err != nil {
                return err
            }
            if err := os.Chmod(destPath, sfi.Mode()); err != nil {
                return err
            }
        }

        return nil
    })
    if err != nil {
        return
    }

    return
}
