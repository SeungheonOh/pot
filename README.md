# Pot ( Pixel On Terminal)
**These are Pixels, but in Terminals**
Light terminal media viewer with modularized loader and render engine. 

## Requirement
Compile time
```
Nothing but standard Go library
```
Runtime 
```
ffmpeg (for ffmpeg based loader)
```

## Modularity
Upon many re-desiging and planning, Pot got modular desigin where users can add their own media loader 
and render engine. This modularity gives ease of maintenance, simpler expension of features and customizability.

### Loaders
PixelOnTerminal only has one loader, ffmpeg based loader. However additional loaders, like linux pipe 
based loader, can be easily depolyed.

The FFmpeg based loader is cabable of loading codecs and formats supported by ffmpeg--```ffmpeg -codecs & ffmpeg -formats```

Currently, only buffered loaders are supported, but non-buffered loader is on its way for live feeds.

### Render Engines
It have got even easier to make a new render engine because now it works upon Go ```image.Image```!
The universial sampling support is comming which will make adding new glyphs on the render engine easier.

## Examples

![Earth](https://github.com/SeungheonOh/pot/blob/master/doc/earth.png)
This image was rendered by 4x8 sampling render engine.

## Things to do (Maybe you can contribute)
- [x] Proper CLI system. 
- [x] Unified input selection for video, image, and cam.
- [ ] New loaders for video live feed.
- [ ] Universial sampling support.
- [ ] Get Some Stars (Yes press that Star NOW)

## Demo
[Demo at Reddit](https://www.reddit.com/r/unixporn/comments/d1gksi/oc_fully_terminal_based_webcamvideoimage_viewer/?utm_source=share&utm_medium=web2x)
This demo is one of the earliest version of PixelOnTerminal and is outdated.

## Inspired by
[Pastel](https://github.com/sharkdp/pastel) -> Definatly one of my favorate CLI app!

## Thanks to
[Diamondburned](https://github.com/diamondburned)
