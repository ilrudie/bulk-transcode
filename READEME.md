# README

## Why?

This project exists to assist with bulk transcoding. Rather than a click-ops setup with handbrake or another GUI, you can use this to help with bulk ffmpeg tasks.

## TODO

- Don't hardcode `mp4` as the only extension we care about.
- Actually execute the commands in some sort of controlled manner, handling non-success results in some sort of reasonable way.
- Watch mode, setup some sort of watching or periodic scanning of the input directory which will transcode any new files.

## Defaults

The defaults baked into the project are for me, on my AMD system. They will transcode video into hevc (H265) without altering the resolution. They simply copy audio.

## Running the project

```shell
❯ go run ./src --input-dir testdata --output-dir testdata/out -r
time=2025-10-26T15:18:42.797-04:00 level=WARN msg="No config file specified, using default configuration"
time=2025-10-26T15:18:42.798-04:00 level=WARN msg="Skipping output directory during recursive scan" output_dir=testdata/out
----------------------------------------
ffmpeg -vaapi_device /dev/dri/renderD128 -hwaccel vaapi -vf 'format=nv12,hwupload' -c:v hevc_vaapi -qp 28 -c:a copy -i testdata/video.mp4 testdata/out/video.hevc.mp4
ffmpeg -vaapi_device /dev/dri/renderD128 -hwaccel vaapi -vf 'format=nv12,hwupload' -c:v hevc_vaapi -qp 28 -c:a copy -i testdata/video1.mp4 testdata/out/video1.hevc.mp4
ffmpeg -vaapi_device /dev/dri/renderD128 -hwaccel vaapi -vf 'format=nv12,hwupload' -c:v hevc_vaapi -qp 28 -c:a copy -i testdata/video2.mp4 testdata/out/video2.hevc.mp4
----------------------------------------
```