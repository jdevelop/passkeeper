EESchema Schematic File Version 2
LIBS:power
LIBS:device
LIBS:transistors
LIBS:conn
LIBS:linear
LIBS:regul
LIBS:74xx
LIBS:cmos4000
LIBS:adc-dac
LIBS:memory
LIBS:xilinx
LIBS:microcontrollers
LIBS:dsp
LIBS:microchip
LIBS:analog_switches
LIBS:motorola
LIBS:texas
LIBS:intel
LIBS:audio
LIBS:interface
LIBS:digital-audio
LIBS:philips
LIBS:display
LIBS:cypress
LIBS:siliconi
LIBS:opto
LIBS:atmel
LIBS:contrib
LIBS:valves
LIBS:test-cad-cache
EELAYER 25 0
EELAYER END
$Descr A4 11693 8268
encoding utf-8
Sheet 1 1
Title ""
Date "2017-05-06"
Rev ""
Comp ""
Comment1 ""
Comment2 ""
Comment3 ""
Comment4 ""
$EndDescr
$Comp
L LCD16X2 DS?
U 1 1 590E7D4B
P 2650 1100
F 0 "DS?" H 1850 1500 50  0000 C CNN
F 1 "LCD16X2" H 3350 1500 50  0000 C CNN
F 2 "WC1602A" H 2650 1050 50  0001 C CIN
F 3 "" H 2650 1100 50  0001 C CNN
	1    2650 1100
	1    0    0    -1  
$EndComp
Wire Wire Line
	3400 1950 3400 1600
Wire Wire Line
	1900 1950 3400 1950
Wire Wire Line
	1900 1600 1900 2100
Connection ~ 1900 1950
Wire Wire Line
	2000 1600 2000 2100
Wire Wire Line
	2200 1600 2200 2100
Wire Wire Line
	2400 1600 2400 2100
Wire Wire Line
	3300 2000 3300 1600
Wire Wire Line
	1600 2000 3300 2000
Connection ~ 2000 2000
Wire Wire Line
	3200 1600 3200 2100
Wire Wire Line
	3200 2100 2900 2100
Wire Wire Line
	3100 1600 3100 2050
Wire Wire Line
	3100 2050 2800 2050
Wire Wire Line
	2800 2050 2800 2100
Wire Wire Line
	3000 1600 3000 1900
Wire Wire Line
	3000 1900 2700 1900
Wire Wire Line
	2700 1900 2700 2100
Wire Wire Line
	2900 1600 2900 1850
Wire Wire Line
	2900 1850 2600 1850
Wire Wire Line
	2600 1850 2600 2100
$Comp
L R 2.2KOm
U 1 1 5910552A
P 1600 1850
F 0 "2.2KOm" V 1680 1850 50  0000 C CNN
F 1 "R" V 1600 1850 50  0000 C CNN
F 2 "" V 1530 1850 50  0001 C CNN
F 3 "" H 1600 1850 50  0001 C CNN
	1    1600 1850
	1    0    0    -1  
$EndComp
Wire Wire Line
	2100 1600 2100 1700
Wire Wire Line
	2100 1700 1600 1700
$Comp
L LED D?
U 1 1 59152D0F
P 1700 2700
F 0 "D?" H 1700 2800 50  0000 C CNN
F 1 "LED" H 1700 2600 50  0000 C CNN
F 2 "" H 1700 2700 50  0001 C CNN
F 3 "" H 1700 2700 50  0001 C CNN
	1    1700 2700
	1    0    0    -1  
$EndComp
$Comp
L LED D?
U 1 1 59152D62
P 2300 2700
F 0 "D?" H 2300 2800 50  0000 C CNN
F 1 "LED" H 2300 2600 50  0000 C CNN
F 2 "" H 2300 2700 50  0001 C CNN
F 3 "" H 2300 2700 50  0001 C CNN
	1    2300 2700
	1    0    0    -1  
$EndComp
$Comp
L LED D?
U 1 1 59152D9B
P 2850 2700
F 0 "D?" H 2850 2800 50  0000 C CNN
F 1 "LED" H 2850 2600 50  0000 C CNN
F 2 "" H 2850 2700 50  0001 C CNN
F 3 "" H 2850 2700 50  0001 C CNN
	1    2850 2700
	1    0    0    -1  
$EndComp
$Comp
L LED D?
U 1 1 59152DC0
P 3500 2700
F 0 "D?" H 3500 2800 50  0000 C CNN
F 1 "LED" H 3500 2600 50  0000 C CNN
F 2 "" H 3500 2700 50  0001 C CNN
F 3 "" H 3500 2700 50  0001 C CNN
	1    3500 2700
	1    0    0    -1  
$EndComp
$Comp
L GND #PWR?
U 1 1 59152E59
P 1150 2600
F 0 "#PWR?" H 1150 2350 50  0001 C CNN
F 1 "GND" H 1150 2450 50  0000 C CNN
F 2 "" H 1150 2600 50  0001 C CNN
F 3 "" H 1150 2600 50  0001 C CNN
	1    1150 2600
	1    0    0    -1  
$EndComp
Wire Wire Line
	1550 2700 1150 2700
Wire Wire Line
	1150 2600 1150 2900
Wire Wire Line
	2150 2700 2050 2700
Wire Wire Line
	2050 2700 2050 2900
Wire Wire Line
	1150 2900 3350 2900
Wire Wire Line
	2700 2900 2700 2700
Connection ~ 2050 2900
Wire Wire Line
	3350 2900 3350 2700
Connection ~ 2700 2900
Wire Wire Line
	1900 2100 1150 2100
Wire Wire Line
	1150 2100 1150 2650
Connection ~ 1150 2650
Wire Wire Line
	2300 1600 2300 1950
Connection ~ 2300 1950
$EndSCHEMATC