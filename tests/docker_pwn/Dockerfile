FROM ubuntu:18.04

RUN echo ' deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse \n \
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse \n \ 
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse \n \
deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse \n '> /etc/apt/sources.list

RUN apt update
RUN apt install -y xinetd

RUN useradd -m ctf

COPY ./flag /flag
COPY ./pwn /pwn
COPY ./ctf.xinetd /etc/xinetd.d/ctf
COPY ./changeFlag.sh /changeFlag.sh
RUN chown root:ctf /pwn /changeFlag.sh && chmod 751 /pwn /changeFlag.sh 

RUN echo 'ctf - nproc 1500' >>/etc/security/limits.conf

CMD exec /bin/bash -c "/etc/init.d/xinetd start; trap : TERM INT; sleep infinity & wait"

EXPOSE 2817