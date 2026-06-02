FROM ubuntu:22.04
RUN useradd -m server
WORKDIR /home/server
USER server
COPY --chown=server LinuxServer .
RUN chmod +x TagDuelsServer.sh
CMD ./TagDuelsServer.sh -ip=0.0.0.0 -ini:Engine:[OnlineSubsystem]:DefaultPlatformService=Null -ini:Game:[OnlineSubsystemAddresses]:UserService=147.45.183.147:8080  -ini:Game:[OnlineSubsystemAddresses]:UserEvent=5.23.53.58:8080 -ini:Game:[OnlineSubsystemAddresses]:MatchHistoryService=82.97.240.43:8080 -ini:Game:[OnlineSubsystemAddresses]:MMEvent=213.171.7.70:8080 -ini:Game:[OnlineSubsystemAddresses]:RatingService=37.220.80.247:8080 -ini:Game:[OnlineSubsystemAddresses]:MatchService=83.147.246.152:8080
