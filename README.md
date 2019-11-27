# DIY Hardware Password Keeper

## Contributors
* [Eugene](https://github.com/jdevelop)
* [Maxim](https://github.com/maximcamario)

## Objective
The goal of the project is to create the easy to build and use device that allows to store virtually unlimited number of passwords on AES encrypted storage. The storage is encrypted with AES and uses the password that is stored on the RFID FOB key.

## Hardware
Hardware part is built on [Raspberry Pi Zero](https://www.raspberrypi.org/products/raspberry-pi-zero/) ([W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) - optional). 

### Bill of material

* [Raspberry Pi Zero](https://www.raspberrypi.org/products/raspberry-pi-zero/) ([W](https://www.raspberrypi.org/products/raspberry-pi-zero-w/) - optional)
* [OLED Display - 0.96" I2C IIC SPI Serial 128X64 OLED LCD Display SSD1306 for 51 STM32 Arduino](https://www.ebay.com/itm/0-96-I2C-IIC-SPI-Serial-128X64-OLED-LCD-Display-SSD1306-for-51-STM32-Arduino-/201688735605)
*  [MFRC-522 RC522 RFID Radiofrequency IC Card Inducing Sensor Reader for Arduino M5](https://www.ebay.com/itm/MFRC-522-RC522-RFID-Radiofrequency-IC-Card-Inducing-Sensor-Reader-for-Arduino-M5-/301723476083)
* [Momentary Panel PCB Tactile Tact Push Button Switch 4Pin 6x6x4.5mm TS](https://www.ebay.com/itm/100Pcs-Momentary-Panel-PCB-Tactile-Tact-Push-Button-Switch-4Pin-6x6x4-5mm-TS-/172040053810?hash=item280e62e432)
* [3mm LED Light White Yellow Red Green Assortment Kit for Arduino](https://www.ebay.com/itm/3mm-5mm-Assortment-LED-Diodes-Light-Emitting-Kit-Red-Green-Blue-Yellow-White-/153414025406)
* [0.8mm 63/37 Tin Lead Rosin Core Solder Flux Soldering Welding Iron Wire Reel 14m](https://www.ebay.com/itm/0-8mm-63-37-Tin-Lead-Rosin-Core-Solder-Flux-Soldering-Welding-Iron-Wire-Reel-14m-/172519124561)
* [AWG30 wrapping wire](https://www.ebay.com/itm/AWG30-Insulated-Wire-Wrapping-Wires-Reel-250M-White-P-N-B-30-1000-K4F8/263834804615)
* [4 x M2x12 bolts/nuts](https://www.ebay.com/itm/480pcs-M2-M3-M4-Metric-Hex-Socket-Head-Cap-Screw-Bolts-Nuts-Assorted-Box-Kit-Set-/152695878844)
* Soldering iron

## Firmware

The OS is based on [Linux Kernel 4.14](https://github.com/torvalds/linux/commit/865ddc1393f558198e7e7ce70928ff2e49c4f7f6) for Raspberry Pi. 

The application is written in [Golang](https://golang.org/) with heavy use of [https://periph.io/](https://periph.io/) components for MFRC522 and OLED display.

### Building

#### Requirements

* [GNU Make](https://www.gnu.org/software/make/)
* [Docker](https://www.docker.com/)

#### Instructions

The application consists of 2 components:

* firmware 
    * service
    * oled splash
* web interface ( AngularJS )


```
make all
``` 
will build the web interface, services and will invoke buildroot to produce the disk image in `dist/sdcard.img.gz` file. This file can be written onto a micro SD card as `dd if=dist/sdcard.img of=/dev/sda bs=1M`

# Testing

Plug PassKeeper into any USB port of the computer. PassKeeper will register itself and will be accessible at http://10.101.1.1/