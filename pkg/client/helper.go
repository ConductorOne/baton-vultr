package client

import (
	"context"
	"golang.org/x/oauth2"
	"log"
	"net/url"
	"strconv"
)

const ItemsPerPage = 5

type PageOptions struct {
	PageSize  int
	PageToken string
}

func getTokenSource(ctx context.Context, bearerToken string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: bearerToken,
		TokenType:   "Bearer",
	})
}

func getPageSize(pageSize int) int {
	if pageSize <= 0 || pageSize > ItemsPerPage {
		pageSize = ItemsPerPage
	}
	return pageSize
}

func getNextPageToken(prevPageToken string, pageSize, totalRecords, recordsCount int) string {
	if prevPageToken == "" {
		return "0"
	}

	prevToken, err := strconv.Atoi(prevPageToken)
	if err != nil {
		log.Printf("Error al convertir prevPageToken a entero: %v", err)
		return ""
	}

	pageSize = getPageSize(pageSize)

	if recordsCount < pageSize || prevToken+pageSize >= totalRecords {
		return "" // No more pages
	}

	nextToken := strconv.Itoa(prevToken + pageSize)
	return nextToken
}

func WithPageLimit(pageSize int) ReqOpt {
	return WithQueryParam("limit", strconv.Itoa(getPageSize(pageSize)))
}

func WithPageIndex(nextPageToken string) ReqOpt {
	if nextPageToken == "" {
		nextPageToken = "0"
	}
	return WithQueryParam("sIndex", nextPageToken)
}

func WithQueryParam(key string, value string) ReqOpt {
	return func(reqURL *url.URL) {
		q := reqURL.Query()
		q.Set(key, value)
		reqURL.RawQuery = q.Encode()
	}
}

type ReqOpt func(reqURL *url.URL)
