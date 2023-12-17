FROM scratch
ADD fruits_api /app/fruits_api
CMD [ "/app/fruits_api" ]