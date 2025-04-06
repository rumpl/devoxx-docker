# Outro

You made it! Pat yourself on the back, you started from nothing and now can run
something that looks like a container. And hopefully you've learned something.

# Next steps

The learning should never stop, here are some extra things you could do with the
code you have right now. Go nuts!

## Layers

In the workshop we used the complete rootfs of an image that was layed out for
us thanks to the `remote` package. If you look at the `/fs/alpine` directory,
you will find multiple things

```console
$ ls -l /fs/alpine/
total 12
-rw-r--r--  1 root root  597 Apr  6 10:30 config.json
drwxr-xr-x 19 root root 4096 Apr  6 10:30 rootfs
drwxr-xr-x  2 root root 4096 Apr  6 10:30 sha256:6e771e15690e2fabf2332d3a3b744495411d6e0b00b2aea64419b58b0066cf81
```

We can see:

- `rootfs` is the extracted root filesystem of the image
- `config.json` is the OCI configuration of the image
- `sha:...` is one layer of the alpine image

The `config.json` file, among other things, defines the layers to use to make
the rootfs of the image, for `alpine` this is:
```json
"rootfs": {
    "type": "layers",
    "diff_ids": [
        "sha256:a16e98724c05975ee8c40d8fe389c3481373d34ab20a1cf52ea2accc43f71f4c"
    ]
},
```

Now try to change your runtime so that you use these layers, each extracted in
its own directory.

<details>
<summary>Hint</summary>

The `diff_ids` is a list of digest of the _uncompressed_ data in the layer but
the digest of a layer is the digest of the _compressed_ data. When
uncompressing, make sure to calculate the uncompressed data digest so that you
can find the layers with the `diff_ids`.

</details>

## More OCI

Browse around the `config.json` file, use parts of it, any part you want,
`config.Cmd` is a nice start isn't it?

## Even more OCI

Read the OCI spec(s), pick something, play with it

## Even more

With the knowledge you have gained now, go ahead and clone any container runtime
out there, read its code, maybe contribute? ;)

- [youki](https://github.com/youki-dev/youki)
- [runc](https://github.com/opencontainers/runc)
- [crun](https://github.com/containers/crun)

And there are many more!