FROM ubuntu:22.04
RUN useradd -m server
WORKDIR /home/server
USER server
COPY LinuxServer .
RUN chmod +x TagDuelsServer.sh
CMD ./TagDuelsServer.sh -ip=0.0.0.0 -ini:Engine:[OnlineSubsystem]:DefaultPlatformService=Null
