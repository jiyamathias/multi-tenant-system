// Package helper contains utility functions and variables that can be used in other packages
package helper

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"codematic/model/pagination"
)

type (
	// SortOrder struct
	SortOrder string
	// Key is a middleware key sting value
	Key string
)

const (
	// ZeroUUID default empty or non set UUID value
	ZeroUUID = "00000000-0000-0000-0000-000000000000"
	// LogStrKeyModule log service name value
	LogStrKeyModule = "ser_name"
	// LogStrKeyLevel log service level value
	LogStrKeyLevel = "lev_name"
	// LogStrPartnerLevel log partner name value
	LogStrPartnerLevel = "partner_name"
	// LogStrKeyMethod log method name value
	LogStrKeyMethod = "method_name"
	// LogStrPackageLevel log package name value
	LogStrPackageLevel = "package_name"
	// LogStrKeyEndpointName log endpoint name value
	LogStrKeyEndpointName = "endpoint_name"
	// SortOrderASC for ascending sorting
	SortOrderASC SortOrder = "ASC"
	// SortOrderDESC for descending sorting
	SortOrderDESC SortOrder = "DESC"
	// GinContextKey constant that holds the Gin context key
	GinContextKey Key = "CodeMatic_GinContextKey"
)

// GetStringPointer returns a string pointer
func GetStringPointer(val string) *string {
	return &val
}

// GetTimePointer returns a time pointer
func GetTimePointer(time time.Time) *time.Time {
	return &time
}

// GetStringVal return string from pointer
func GetStringVal(strVal *string) string {
	var val string
	if strVal != nil {
		return *strVal
	}
	return val
}

// GetIntVal return int from pointer
func GetIntVal(intVal *int) int {
	var val int
	if intVal != nil {
		return *intVal
	}
	return val
}

// GetFloat64Pointer returns a float64 pointer
func GetFloat64Pointer(val float64) *float64 {
	return &val
}

// GetFloatVal returns valid float val from pointer
func GetFloatVal(floatVal *float64) float64 {
	var val float64
	if floatVal != nil {
		return *floatVal
	}
	return val
}

// GetBoolPointer returns a bool pointer
func GetBoolPointer(val bool) *bool {
	return &val
}

// GetBoolVal returns a bool value from pointer
func GetBoolVal(boolVal *bool) bool {
	var val bool
	if boolVal != nil {
		return *boolVal
	}
	return val
}

// GetIntPointer returns an int pointer
func GetIntPointer(val int) *int {
	return &val
}

// RandomNumbers generates random numerics with length specified
func RandomNumbers(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b)
}

// GenerateRandomString generates a random string
func GenerateRandomString(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		randomString[i] = charset[randomIndex.Int64()]
	}

	return string(randomString), nil
}

// GenerateKey generate random key
func GenerateKey(max int) (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	b := make([]byte, max)
	_, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

// GenerateRandomDigits generate random digits
func GenerateRandomDigits(max int) (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	b := make([]byte, max)
	_, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

// ParsePageParams get the pagination details from the gin context
func ParsePageParams(c *gin.Context) pagination.Page {
	size := pagination.PageDefaultSize
	pageNo := pagination.PageDefaultNumber
	sortDirection := pagination.PageDefaultSortDirectionDesc
	sortBy := pagination.PageDefaultSortBy

	s, err := strconv.Atoi(c.Query("size"))
	if err == nil {
		size = s
	}

	n, err := strconv.Atoi(c.Query("page"))
	if err == nil {
		pageNo = n
	}

	sD, err := strconv.ParseBool(c.Query("sort_direction_desc"))
	if err == nil {
		sortDirection = sD
	}

	sB := c.Query("sort_by")
	if !strings.EqualFold(sB, "") {
		sortBy = sB
	}

	return pagination.Page{
		Number:            &pageNo,
		Size:              &size,
		SortBy:            &sortBy,
		SortDirectionDesc: &sortDirection,
	}
}

// CustomError sends a custom error
func CustomError(err string) error {
	if len(err) == 0 {
		return ErrDefault
	}

	return errors.New(err)
}
