FROM centos:7

EXPOSE 9090

ADD ./bin/woocommerce-customer-service /usr/bin/woocommerce-customer-service

CMD ["woocommerce-customer-service", "serve"]