// Code generated by fileb0x at "2018-12-12 17:37:49.685 +0900 KST m=+0.207000001" from config file "aScript.json" DO NOT EDIT.
// modification hash(684420c151919b6365008860dcec43bb.55d3f29a56c9edb41eb9ef813100e6a0)

package shader

import (
	"bytes"

	"context"
	"io"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/webdav"
)

var (
	// CTX is a context for webdav vfs
	CTX = context.Background()

	// FS is a virtual memory file system
	FS = webdav.NewMemFS()

	// Handler is used to server files through a http handler
	Handler *webdav.Handler

	// HTTP is the http file system
	HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct {
	// Prefix allows to limit the path of all requests. F.e. a prefix "css" would allow only calls to /css/*
	Prefix string
}

// FileDebugFsGlsl is "Debug.fs.glsl"
var FileDebugFsGlsl = []byte("\x0d\x0a\x76\x6f\x69\x64\x20\x6d\x61\x69\x6e\x28\x29\x20\x7b\x0d\x0a\x0d\x0a\x7d\x0d\x0a")

// FileFlatFsGlsl is "Flat.fs.glsl"
var FileFlatFsGlsl = []byte("\x0d\x0a\x75\x6e\x69\x66\x6f\x72\x6d\x20\x76\x65\x63\x34\x20\x42\x61\x73\x65\x43\x6f\x6c\x6f\x72\x46\x61\x63\x74\x6f\x72\x3b\x0d\x0a\x0d\x0a\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x42\x41\x53\x45\x43\x4f\x4c\x4f\x52\x54\x45\x58\x0d\x0a\x75\x6e\x69\x66\x6f\x72\x6d\x20\x73\x61\x6d\x70\x6c\x65\x72\x32\x44\x20\x42\x61\x73\x65\x43\x6f\x6c\x6f\x72\x54\x65\x78\x3b\x0d\x0a\x23\x65\x6e\x64\x69\x66\x0d\x0a\x0d\x0a\x69\x6e\x20\x73\x74\x72\x75\x63\x74\x7b\x0d\x0a\x20\x20\x20\x76\x65\x63\x33\x20\x70\x6f\x73\x69\x74\x69\x6f\x6e\x3b\x0d\x0a\x20\x20\x20\x76\x65\x63\x32\x20\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x3b\x0d\x0a\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x4e\x4f\x52\x4d\x41\x4c\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x54\x41\x4e\x47\x45\x4e\x54\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x6d\x61\x74\x33\x20\x54\x42\x4e\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x23\x65\x6c\x73\x65\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x76\x65\x63\x33\x20\x6e\x6f\x72\x6d\x61\x6c\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x7d\x20\x66\x73\x6f\x75\x74\x3b\x0d\x0a\x0d\x0a\x0d\x0a\x0d\x0a\x6f\x75\x74\x20\x76\x65\x63\x34\x20\x6f\x75\x74\x70\x75\x74\x43\x6f\x6c\x6f\x72\x3b\x0d\x0a\x0d\x0a\x76\x6f\x69\x64\x20\x6d\x61\x69\x6e\x28\x29\x20\x7b\x0d\x0a\x20\x20\x20\x20\x76\x65\x63\x34\x20\x63\x20\x3d\x20\x42\x61\x73\x65\x43\x6f\x6c\x6f\x72\x46\x61\x63\x74\x6f\x72\x3b\x0d\x0a\x0d\x0a\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x42\x41\x53\x45\x43\x4f\x4c\x4f\x52\x54\x45\x58\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x63\x20\x3d\x20\x63\x20\x2a\x20\x74\x65\x78\x74\x75\x72\x65\x28\x73\x61\x6d\x70\x6c\x65\x72\x32\x44\x28\x30\x29\x2c\x20\x66\x73\x6f\x75\x74\x2e\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x2e\x78\x79\x29\x3b\x0d\x0a\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x20\x20\x20\x20\x6f\x75\x74\x70\x75\x74\x43\x6f\x6c\x6f\x72\x20\x3d\x20\x63\x3b\x0d\x0a\x7d\x0d\x0a")

// FilePBRFsGlsl is "PBR.fs.glsl"
var FilePBRFsGlsl = []byte("\x0d\x0a\x76\x6f\x69\x64\x20\x6d\x61\x69\x6e\x28\x29\x20\x7b\x0d\x0a\x0d\x0a\x7d\x0d\x0a")

// FilePhongFsGlsl is "Phong.fs.glsl"
var FilePhongFsGlsl = []byte("\x0d\x0a\x0d\x0a\x76\x6f\x69\x64\x20\x6d\x61\x69\x6e\x28\x29\x20\x7b\x0d\x0a\x0d\x0a\x7d\x0d\x0a")

// FileStandardVsGlsl is "Standard.vs.glsl"
var FileStandardVsGlsl = []byte("\x0d\x0a\x0d\x0a\x2f\x2f\x0d\x0a\x2f\x2f\x20\x52\x65\x66\x65\x72\x65\x6e\x63\x65\x20\x3e\x20\x68\x74\x74\x70\x73\x3a\x2f\x2f\x67\x69\x74\x68\x75\x62\x2e\x63\x6f\x6d\x2f\x4b\x68\x72\x6f\x6e\x6f\x73\x47\x72\x6f\x75\x70\x2f\x67\x6c\x54\x46\x2d\x57\x65\x62\x47\x4c\x2d\x50\x42\x52\x2f\x62\x6c\x6f\x62\x2f\x6d\x61\x73\x74\x65\x72\x2f\x73\x68\x61\x64\x65\x72\x73\x2f\x70\x62\x72\x2d\x76\x65\x72\x74\x2e\x67\x6c\x73\x6c\x0d\x0a\x2f\x2f\x0d\x0a\x0d\x0a\x75\x6e\x69\x66\x6f\x72\x6d\x20\x6d\x61\x74\x34\x20\x43\x61\x6d\x65\x72\x61\x4d\x61\x74\x72\x69\x78\x3b\x0d\x0a\x75\x6e\x69\x66\x6f\x72\x6d\x20\x6d\x61\x74\x34\x20\x4d\x6f\x64\x65\x6c\x4d\x61\x74\x72\x69\x78\x3b\x0d\x0a\x75\x6e\x69\x66\x6f\x72\x6d\x20\x6d\x61\x74\x34\x20\x4e\x6f\x72\x6d\x61\x6c\x4d\x61\x74\x72\x69\x78\x3b\x0d\x0a\x0d\x0a\x6c\x61\x79\x6f\x75\x74\x20\x28\x6c\x6f\x63\x61\x74\x69\x6f\x6e\x20\x3d\x20\x30\x29\x20\x69\x6e\x20\x76\x65\x63\x33\x20\x70\x6f\x73\x69\x74\x69\x6f\x6e\x3b\x0d\x0a\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x4e\x4f\x52\x4d\x41\x4c\x0d\x0a\x6c\x61\x79\x6f\x75\x74\x20\x28\x6c\x6f\x63\x61\x74\x69\x6f\x6e\x20\x3d\x20\x31\x29\x20\x69\x6e\x20\x76\x65\x63\x33\x20\x6e\x6f\x72\x6d\x61\x6c\x3b\x0d\x0a\x23\x65\x6e\x64\x69\x66\x0d\x0a\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x54\x41\x4e\x47\x45\x4e\x54\x0d\x0a\x6c\x61\x79\x6f\x75\x74\x20\x28\x6c\x6f\x63\x61\x74\x69\x6f\x6e\x20\x3d\x20\x32\x29\x20\x69\x6e\x20\x76\x65\x63\x33\x20\x74\x61\x6e\x67\x65\x6e\x74\x3b\x0d\x0a\x23\x65\x6e\x64\x69\x66\x0d\x0a\x2f\x2f\x0d\x0a\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x43\x4f\x4f\x52\x44\x5f\x30\x0d\x0a\x6c\x61\x79\x6f\x75\x74\x20\x28\x6c\x6f\x63\x61\x74\x69\x6f\x6e\x20\x3d\x20\x34\x29\x20\x69\x6e\x20\x76\x65\x63\x32\x20\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x3b\x0d\x0a\x23\x65\x6e\x64\x69\x66\x0d\x0a\x2f\x2f\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x43\x4f\x4f\x52\x44\x5f\x31\x0d\x0a\x2f\x2f\x6c\x61\x79\x6f\x75\x74\x20\x28\x6c\x6f\x63\x61\x74\x69\x6f\x6e\x20\x3d\x20\x35\x29\x20\x69\x6e\x20\x76\x65\x63\x32\x20\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x31\x3b\x0d\x0a\x2f\x2f\x23\x65\x6e\x64\x69\x66\x0d\x0a\x0d\x0a\x0d\x0a\x6f\x75\x74\x20\x73\x74\x72\x75\x63\x74\x7b\x0d\x0a\x20\x20\x20\x20\x76\x65\x63\x33\x20\x70\x6f\x73\x69\x74\x69\x6f\x6e\x3b\x0d\x0a\x20\x20\x20\x20\x76\x65\x63\x32\x20\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x3b\x0d\x0a\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x4e\x4f\x52\x4d\x41\x4c\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x54\x41\x4e\x47\x45\x4e\x54\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x6d\x61\x74\x33\x20\x54\x42\x4e\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x23\x65\x6c\x73\x65\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x76\x65\x63\x33\x20\x6e\x6f\x72\x6d\x61\x6c\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x7d\x20\x66\x73\x6f\x75\x74\x3b\x0d\x0a\x0d\x0a\x76\x6f\x69\x64\x20\x6d\x61\x69\x6e\x28\x29\x20\x7b\x0d\x0a\x20\x20\x20\x20\x76\x65\x63\x34\x20\x70\x6f\x73\x20\x3d\x20\x4d\x6f\x64\x65\x6c\x4d\x61\x74\x72\x69\x78\x20\x2a\x20\x76\x65\x63\x34\x28\x70\x6f\x73\x69\x74\x69\x6f\x6e\x2c\x20\x31\x29\x3b\x0d\x0a\x20\x20\x20\x20\x66\x73\x6f\x75\x74\x2e\x70\x6f\x73\x69\x74\x69\x6f\x6e\x20\x3d\x20\x76\x65\x63\x33\x28\x70\x6f\x73\x2e\x78\x79\x7a\x29\x20\x2f\x20\x70\x6f\x73\x2e\x77\x3b\x0d\x0a\x20\x20\x20\x20\x2f\x2f\x20\x4e\x6f\x72\x6d\x61\x6c\x2c\x0d\x0a\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x4e\x4f\x52\x4d\x41\x4c\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x54\x41\x4e\x47\x45\x4e\x54\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x2f\x2f\x20\x69\x66\x20\x48\x41\x53\x5f\x4e\x4f\x52\x4d\x41\x4c\x20\x61\x6e\x64\x20\x48\x41\x53\x5f\x54\x41\x4e\x47\x45\x4e\x54\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x76\x65\x63\x33\x20\x6e\x6f\x72\x6d\x61\x6c\x57\x20\x3d\x20\x6e\x6f\x72\x6d\x61\x6c\x69\x7a\x65\x28\x76\x65\x63\x33\x28\x4e\x6f\x72\x6d\x61\x6c\x4d\x61\x74\x72\x69\x78\x20\x2a\x20\x76\x65\x63\x34\x28\x6e\x6f\x72\x6d\x61\x6c\x2e\x78\x79\x7a\x2c\x20\x30\x2e\x30\x29\x29\x29\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x76\x65\x63\x33\x20\x74\x61\x6e\x67\x65\x6e\x74\x57\x20\x3d\x20\x6e\x6f\x72\x6d\x61\x6c\x69\x7a\x65\x28\x76\x65\x63\x33\x28\x4d\x6f\x64\x65\x6c\x4d\x61\x74\x72\x69\x78\x20\x2a\x20\x76\x65\x63\x34\x28\x74\x61\x6e\x67\x65\x6e\x74\x2e\x78\x79\x7a\x2c\x20\x30\x2e\x30\x29\x29\x29\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x76\x65\x63\x33\x20\x62\x69\x74\x61\x6e\x67\x65\x6e\x74\x57\x20\x3d\x20\x63\x72\x6f\x73\x73\x28\x6e\x6f\x72\x6d\x61\x6c\x57\x2c\x20\x74\x61\x6e\x67\x65\x6e\x74\x57\x29\x20\x2a\x20\x74\x61\x6e\x67\x65\x6e\x74\x2e\x77\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x66\x73\x6f\x75\x74\x2e\x54\x42\x4e\x20\x3d\x20\x6d\x61\x74\x33\x28\x74\x61\x6e\x67\x65\x6e\x74\x57\x2c\x20\x62\x69\x74\x61\x6e\x67\x65\x6e\x74\x57\x2c\x20\x6e\x6f\x72\x6d\x61\x6c\x57\x29\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x23\x65\x6c\x73\x65\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x2f\x2f\x20\x69\x66\x20\x48\x41\x53\x5f\x4e\x4f\x52\x4d\x41\x4c\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x66\x73\x6f\x75\x74\x2e\x6e\x6f\x72\x6d\x61\x6c\x20\x3d\x20\x6e\x6f\x72\x6d\x61\x6c\x69\x7a\x65\x28\x76\x65\x63\x33\x28\x4d\x6f\x64\x65\x6c\x4d\x61\x74\x72\x69\x78\x20\x2a\x20\x76\x65\x63\x34\x28\x6e\x6f\x72\x6d\x61\x6c\x2e\x78\x79\x7a\x2c\x20\x30\x2e\x30\x29\x29\x29\x3b\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x20\x20\x20\x20\x2f\x2f\x20\x54\x65\x78\x43\x6f\x6f\x72\x64\x0d\x0a\x20\x20\x20\x20\x23\x69\x66\x64\x65\x66\x20\x48\x41\x53\x5f\x43\x4f\x4f\x52\x44\x5f\x30\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x66\x73\x6f\x75\x74\x2e\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x20\x3d\x20\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x3b\x0d\x0a\x20\x20\x20\x20\x23\x65\x6c\x73\x65\x0d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x66\x73\x6f\x75\x74\x2e\x74\x65\x78\x43\x6f\x6f\x72\x64\x5f\x30\x20\x3d\x20\x76\x65\x63\x32\x28\x30\x2e\x30\x2c\x30\x2e\x30\x29\x3b\x0d\x0a\x20\x20\x20\x20\x23\x65\x6e\x64\x69\x66\x0d\x0a\x20\x20\x20\x20\x2f\x2f\x20\x43\x61\x6d\x65\x72\x61\x20\x3d\x20\x50\x65\x72\x73\x70\x65\x63\x74\x69\x76\x65\x20\x2a\x20\x56\x69\x65\x77\x0d\x0a\x20\x20\x20\x20\x67\x6c\x5f\x50\x6f\x73\x69\x74\x69\x6f\x6e\x20\x3d\x20\x43\x61\x6d\x65\x72\x61\x4d\x61\x74\x72\x69\x78\x20\x2a\x20\x4d\x6f\x64\x65\x6c\x4d\x61\x74\x72\x69\x78\x20\x2a\x20\x76\x65\x63\x34\x28\x70\x6f\x73\x69\x74\x69\x6f\x6e\x2c\x20\x31\x29\x3b\x20\x2f\x2f\x20\x6e\x65\x65\x64\x73\x20\x77\x20\x66\x6f\x72\x20\x70\x72\x6f\x70\x65\x72\x20\x70\x65\x72\x73\x70\x65\x63\x74\x69\x76\x65\x20\x63\x6f\x72\x72\x65\x63\x74\x69\x6f\x6e\x0d\x0a\x7d")

func init() {
	err := CTX.Err()
	if err != nil {
		panic(err)
	}

	var f webdav.File

	f, err = FS.OpenFile(CTX, "Debug.fs.glsl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FileDebugFsGlsl)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "Flat.fs.glsl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FileFlatFsGlsl)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "PBR.fs.glsl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FilePBRFsGlsl)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "Phong.fs.glsl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FilePhongFsGlsl)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	f, err = FS.OpenFile(CTX, "Standard.vs.glsl", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FileStandardVsGlsl)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	Handler = &webdav.Handler{
		FileSystem: FS,
		LockSystem: webdav.NewMemLS(),
	}

}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {
	path = hfs.Prefix + path

	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(f)
	return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// WalkDirs looks for files in the given dir and returns a list of files in it
// usage for all files in the b0x: WalkDirs("", false)
func WalkDirs(name string, includeDirsInList bool, files ...string) ([]string, error) {
	f, err := FS.OpenFile(CTX, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	fileInfos, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	for _, info := range fileInfos {
		filename := path.Join(name, info.Name())

		if includeDirsInList || !info.IsDir() {
			files = append(files, filename)
		}

		if info.IsDir() {
			files, err = WalkDirs(filename, includeDirsInList, files...)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}
