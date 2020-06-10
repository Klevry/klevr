FROM klevry/klevr-centosbase:latest
MAINTAINER ralfyang <"https://github.com/ralfyang">

LABEL version=0.1
LABEL image_name=klevry/libvirt

RUN yum -y install libvirt-client libvirt-daemon-driver-* libvirt-daemon libvirt-daemon-kvm qemu-kvm && yum clean all
RUN systemctl enable libvirtd; systemctl enable virtlockd && yum install -y openssh-server openssh-clients arpwatch

RUN echo "root:klevr00" |chpasswd
RUN systemctl enable sshd
RUN sed -i 's|[#]*PermitRootLogin no|PermitRootLogin yes|g' /etc/ssh/sshd_config
RUN sed -i 's|[#]*PasswordAuthentication no|PasswordAuthentication yes|g' /etc/ssh/sshd_config
RUN sed -i 's|[#]*ChallengeResponseAuthentication no|ChallengeResponseAuthentication yes|g' /etc/ssh/sshd_config
RUN sed -i 's|UsePAM no|UsePAM yes|g' /etc/ssh/sshd_config

RUN echo "listen_tls = 0" >> /etc/libvirt/libvirtd.conf; \
echo 'listen_tcp = 1' >> /etc/libvirt/libvirtd.conf; \
echo 'tls_port = "16514"' >> /etc/libvirt/libvirtd.conf; \ 
echo 'tcp_port = "16509"' >> /etc/libvirt/libvirtd.conf; \
echo 'auth_tcp = "none"' >> /etc/libvirt/libvirtd.conf

RUN echo 'vnc_listen = "0.0.0.0"' >> /etc/libvirt/qemu.conf
RUN echo 'LIBVIRTD_ARGS="--listen"' >> /etc/sysconfig/libvirtd

COPY network.xml /network.xml
COPY customlibvirtpost.service /usr/lib/systemd/system/customlibvirtpost.service
COPY customlibvirtpost.sh /customlibvirtpost.sh
RUN chmod a+x /customlibvirtpost.sh
RUN systemctl enable customlibvirtpost

EXPOSE 22
EXPOSE 16509
EXPOSE 5900

VOLUME ["/sys/fs/cgroup"]
#CMD  [ "/usr/sbin/init" ]
ENTRYPOINT  ["/usr/sbin/init"]
