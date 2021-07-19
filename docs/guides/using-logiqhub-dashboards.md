# Using LOGIQHub dashboards

LOGIQHub is a collection of pre-built, ready-to-use dashboards for monitoring and observing important metrics for popular infrastructure and cloud services. Using these dashboards, you can go from installing LOGIQ to instantly deriving valuable and actionable insights from logs and metrics from your infrastructure services.

LOGIQ dashboards are preconfigured to extract, plot, and visualize the most important performance, health, and usage metrics from your infrastructure and cloud services so that you don’t have to spend time writing queries and setting them up from scratch. LOGIQHub dashboards are the quickest, no-code approach to getting started with monitoring, observability, and data and metrics visualizations for your applications, services, and infrastructure on LOGIQ.

To use a LOGIQHub dashboard, all you’ll need:

- A LOGIQ instance configured to ingest data from your data sources
- the LOGIQ CLI – logiqctl
- A LOGIQHub dashboard suitable for your service

The following section will help you get started with using a LOGIQHub dashboard for monitoring your tech environment.

## Importing a LOGIQHub dashboard

In order to import a LOGIQHub dashboard into your LOGIQ instance, do the following:

1. Head over to [LOGIQHub](https://logiqhub.logiq.ai/) and access the dashboard of your choice. 
2. Download the dashboard JSON to your machine. 
3. Open the JSON file in a text editor and edit the `datasources` section to add your Prometheus endpoint. 
    1. For Kubernetes, Memcached, MySQL, Prometheus, and Redis dashboard, edit the `namespaces` you'd like to monitor. 
4. Optionally, rename the dashboard. 
5. Import the dashboard into your LOGIQ instance by running the following command, after replacing `<dashboard-file>` with the name of your dashboard JSON file. 
    ```bash
    logiqctl create dashboard -f <dashboard-file>.json
    ```

When you log into your LOGIQ UI and head over to **Dashboards**, you’ll now see a list of newly-created dashboards containing visualizations for all the critical metrics and data for the stack you configured the dashboard for. 

## Supported services

LOGIQHub currently supports pre-built dashboards for the following services. Click on any of the service names to access the dashboard and view service-specific instructions to configure and import the LOGIQHub dashboard. 

- [Fluent Bit](https://github.com/logiqai/logiqhub/tree/master/fluent-bit)
- [Kafka](https://github.com/logiqai/logiqhub/tree/master/kafka)
- [Kubernetes](https://github.com/logiqai/logiqhub/tree/master/kubernetes)
- [LOGIQ](https://github.com/logiqai/logiqhub/tree/master/logiq)
- [Memcached](https://github.com/logiqai/logiqhub/tree/master/memcached)
- [MongoDB](https://github.com/logiqai/logiqhub/tree/master/mongodb)
- [MySQL](https://github.com/logiqai/logiqhub/tree/master/mysql)
- [PostgreSQL](https://github.com/logiqai/logiqhub/tree/master/postgresql)
- [Prometheus](https://github.com/logiqai/logiqhub/tree/master/prometheus)
- [Redis](https://github.com/logiqai/logiqhub/tree/master/redis)

