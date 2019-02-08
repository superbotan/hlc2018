FROM alpine:3.7 
COPY bin/app /opt/api 
COPY bin/run.sh /opt/run.sh 
RUN chmod +x /opt/run.sh
EXPOSE 80
#ENTRYPOINT ["/opt/api"]
ENTRYPOINT ["sh","-c","/opt/run.sh"]

