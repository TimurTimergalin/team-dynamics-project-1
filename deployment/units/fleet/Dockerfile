FROM ubuntu:22.04
RUN useradd -m server
WORKDIR /home/server
COPY LinuxServer .
RUN chown -R server .
RUN chmod +x TagDuelsServer.sh
USER server
CMD ./TagDuelsServer.sh -ip=0.0.0.0 -ini:Engine:[OnlineSubsystem]:DefaultPlatformService=Null
