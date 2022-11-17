This code is intented to run in the background at all times.
It supports a more flexible and stable web conferencing experience.

It opens all /dev/video* devices and looks for one with the "loopback" driver.
It then opens /dev/videoX as a writer and proceeds to write an idle frame every 500ms to it.
This keeps it open and allows Chrome (and other applications?) to find it and use it as a webcam.


------- watch output -------------------------
cvlc v4l2:///dev/videoX

------- v4l2loopback device notes --------
main page: https://github.com/umlaeute/v4l2loopback

Load the loopback device like so:
modprobe v4l2loopback exclusive_caps=yes debug=1

The loopback device has some design considerations that make it
a bit diffult to use for this purpose.
1. It only reports an input as "available" if a writer is currently
writing to it. Annoyingly, that means V4l2 applications looking for a webcam
won't detect it (e.g. Chrome will not display it amongst its camera selection)

2. it has no means (currently) to detect if a reader is present or not.
There is talk about using the events api (https://bugs.launchpad.net/ubuntu/+source/v4l2loopback/+bug/1921474) and some code ( https://github.com/umlaeute/v4l2loopback/blob/main/v4l2loopback.c) - but it is not yet available in my distro (as of 13/11/2022)

------- code hints -------------------------

a loopbackdevice is opened once at the beginning.
it needs two things: a timersource and (at least one) frame provider.
a timersource can be a webcam or an (soft or fake) irq or anything else that provides a regular tick.
a frameprovider provides a frame, one at a time
Whenever the timersource (implemented as a channel) notifies the loopback device it will ask all frameproviders
for their current frames and send them to the output

two frameproviders are currently implemented:
1. IdleFrameProvider (sends an image every 500ms)
2. WebCamFrameProvider (reads an image from a webcam and sends it whenever it got a new one)
