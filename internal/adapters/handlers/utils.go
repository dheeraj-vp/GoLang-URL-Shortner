package handlers

import (
	"strings"

	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/domain"
	"github.com/aws/aws-lambda-go/events"
)

// ExtractPlatformFromRequest determines the platform from request headers
func ExtractPlatformFromRequest(req events.APIGatewayV2HTTPRequest) domain.Platform {
	// Check User-Agent header
	userAgent := strings.ToLower(req.Headers["user-agent"])
	if strings.Contains(userAgent, strings.ToLower(config.InstagramUserAgent)) {
		return domain.PlatformInstagram
	}
	if strings.Contains(userAgent, strings.ToLower(config.TwitterUserAgent)) {
		return domain.PlatformTwitter
	}
	if strings.Contains(userAgent, strings.ToLower(config.YouTubeUserAgent)) {
		return domain.PlatformYouTube
	}

	// Check Referer header
	referer := strings.ToLower(req.Headers["referer"])
	if strings.Contains(referer, config.InstagramReferer) {
		return domain.PlatformInstagram
	}
	if strings.Contains(referer, config.TwitterReferer) {
		return domain.PlatformTwitter
	}
	if strings.Contains(referer, config.YouTubeReferer) {
		return domain.PlatformYouTube
	}

	return domain.PlatformUnknown
}

// IsShortURLLoop checks if the URL is trying to shorten an already shortened URL
func IsShortURLLoop(url, baseURL string) bool {
	return strings.Contains(url, baseURL)
}

// IsMaliciousURL checks for potentially malicious URLs
// This is a basic implementation - in production, you'd want a more sophisticated check
func IsMaliciousURL(url string) bool {
	url = strings.ToLower(url)
	
	// List of suspicious patterns
	suspiciousPatterns := []string{
		"javascript:",
		"data:",
		"file:",
		"vbscript:",
		"about:",
	}
	
	for _, pattern := range suspiciousPatterns {
		if strings.HasPrefix(url, pattern) {
			return true
		}
	}
	
	return false
}
