package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

NSData* getClipboardImageAsPNG() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSArray *classArray = [NSArray arrayWithObject:[NSImage class]];
    NSDictionary *options = [NSDictionary dictionary];

    BOOL exists = [pasteboard canReadObjectForClasses:classArray options:options];
    if (!exists) {
        return nil;
    }

    NSArray *objects = [pasteboard readObjectsForClasses:classArray options:options];
    if (objects == nil || [objects count] == 0) {
        return nil;
    }

    NSImage *image = [objects objectAtIndex:0];
    NSBitmapImageRep *imgRep = [[NSBitmapImageRep alloc] initWithData:[image TIFFRepresentation]];
    NSData *pngData = [imgRep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];

    return pngData;
}

NSUInteger getClipboardImagePNGLength() {
    NSData *data = getClipboardImageAsPNG();
    if (data == nil) {
        return 0;
    }
    return [data length];
}

void getClipboardImagePNGBytes(void *buffer) {
    NSData *data = getClipboardImageAsPNG();
    if (data == nil) {
        return;
    }

    [data getBytes:buffer length:[data length]];
}
*/
import "C"
import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"unsafe"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// GetClipboardImageAsPNG retrieves image content from the macOS clipboard as PNG bytes.
// Returns the PNG data as a byte slice or an error if no image is found.
func GetClipboardImageAsPNG() ([]byte, error) {
	// Check if clipboard has an image
	length := C.getClipboardImagePNGLength()
	if length == 0 {
		return nil, errors.New("no image found in clipboard")
	}

	// Allocate buffer for the PNG data
	buffer := make([]byte, length)

	// Get the PNG data into our buffer
	C.getClipboardImagePNGBytes(unsafe.Pointer(&buffer[0]))

	return buffer, nil
}

func run() error {

	mcpServer := server.NewMCPServer("Snowflake", "1.0.0")

	// Add a query tool.
	mcpServer.AddTool(mcp.NewTool(
		"image_paste",
		mcp.WithDescription("Gets the image from the clipboard."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		b, err := GetClipboardImageAsPNG()
		if err != nil {
			return nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.ImageContent{
					Type:     "image",
					Data:     base64.StdEncoding.EncodeToString(b),
					MIMEType: "image/png",
				},
			},
		}, nil
	})

	return server.ServeStdio(mcpServer)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
