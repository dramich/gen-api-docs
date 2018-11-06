FROM swaggerapi/swagger-ui

ENV SWAGGER_JSON=/data/swagger.json
COPY ./swagger/swagger.out.json /data/swagger.json
