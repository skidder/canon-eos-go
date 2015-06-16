# canon-eos-go
Go bindings for Canon EOS (tm) camera OS X SDK

**This is experimental and under active development, don't expect much**

## Requirements
You must first download the Canon EOS SDK, which involves [requesting access from Canon](http://usa.canon.com/cusa/support/professional/professional_cameras/eos_digital_slr_cameras/eos_7d/standard_display/SDK).

Once downloaded, copy the header and framework files to the system locations:
```shell
# run from the root of the EOS SDK disk image, or where you copied the disk image contents
sudo mkdir /usr/local/include/EDSDK
sudo cp EDSDK/Header/*.h /usr/local/include/EDSDK

sudo cp -r EDSDK/Framework/DPP.framework /Library/Frameworks/
sudo cp -r EDSDK/Framework/EDSDK.framework /Library/Frameworks/
```

## Building
```shell
make all
```
