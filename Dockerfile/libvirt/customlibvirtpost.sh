#/bin/bash

chmod 666 /dev/kvm
virsh net-define /network.xml
virsh net-start klevr
virsh net-autostart klevr
