# eventmaker

Creates and configures Azure IoT Hub with virtual devices and sends mocked events to them. Works both locally, from developer workstation as well as in a fleet from Azure Container Instances (ACI)

![](img/overview.png)

## setup

To run `eventmaker` start by editing a couple variable in the [script/config](script/config) file:

```shell
HUB_NAME="cloudylabs"
NUMBER_OF_INSTANCES=5
METRIC_TEMPLATE="temp|celsius|float|0:72.1|3s"
```

* `HUB_NAME` is the name of the Azure IoT Hub that will be created 
* `NUMBER_OF_INSTANCES` is the number of ACI instances you want to generate. The more instances, the more data will be submitted to Azure IoT Hub
* `METRIC_TEMPLATE` is the template for mocking your data (see section on templates below)

> assumes you have resource group and location defaults set 

Create a IoT Hub with a standard pricing tier (~1 min)

```shell
bin/hub-up
```

## fleet 

Now that the IoT Hub is up, you can deploy the fleet of devices with a corresponding Container Instances in Azure service that send mocked events to these devices. 

```shell
bin/fleet-up
```

> The deployment is asynchronous so if you want to see the result open the ACI dashboard in Azure Portal. Note, may take a ~30 seconds for the first image to appear in the UI

## events 

The mocked event generated by the fleet will look like this

```json
{
    "id": "fdf612b9-34a5-445e-9941-59c404ea9bef",
    "src_id": "client-1",
    "time": 1589745397,
    "label": "temp",
    "data": 70.79129651786347,
    "unit": "celsius"
}
```

## templates

`eventmaker` uses templates to mock the data to submit to Azure IoT Code. For the above event was generated using single metric template like this (`--metric temp|celsius|float|68.9:72.1|3s`). The format of individual metric templates is

`<event label>|<metric unit>|<type of data>|<range of data to generate>|<frequency of sending>`

The supported types are `int`, `float`, and `bool` as well as most common derivates of these (e.g. `int64` or `float32`).

The `ranges` follow `min:max` format. So int in between 0 and 100 would be formatted as `0:100`. This way you can include negative numbers. 

Finally, the `frequency` follows standard go `time.Duration` format (e.g. `1s` for every second, `2m` for every 2 minutes, or `3h` for every 3 hours)

The one defaults you set using environment variables is the device name (`DEV_NAME`) which is the device ID associated with this client (default: `device-1`)

`eventmaker` also supports multiple metrics defined in configuration file. (`--file config/example.yaml`). The format of that file is similar. In the bellow example we have 3 metrics (temp, speed, friction)

```yaml
--- 
metrics: 
- label: temp
  frequency: "3s"
  unit: celsius
  template: 
    max: 107.5
    min: 86.1
    type: float
- label: speed
  frequency: "1s"
  unit: kmh
  template: 
    max: 200
    min: 0
    type: int
- label: friction
  frequency: "1s"
  unit: coefficient
  template: 
    max: 1.00
    min: 0.00
    type: float
```

## endpoints

To find the Azure Service Bus here these events will be published:

```shell
az iot hub show \
  --name $HUB_NAME \
  --query "properties.eventHubEndpoints.events.endpoint" \
  -o tsv
```


## run locally

To run `eventmaker` locally without deployment to ACI

```shell
export HUB_NAME="your-iot-hub-name"
export DIVICE_NAME="device1"
```

Create IoT Hub

```shell
az iot hub create --name $HUB_NAME --sku S1
```

Create the device in the identity registry 

```shell
az iot hub device-identity create \
  --hub-name $HUB_NAME \
  --device-id $DIVICE_NAME
```

Retrieve device connection string

```shell
export CONN_STR=$(az iot hub device-identity show-connection-string \
  --device-id $DIVICE_NAME --hub-name $HUB_NAME -o tsv)
```



If you do make changes to the code you will need to rebuild the executable 

```shell
make build
``` 

To run `eventmaker` locally

```shell
dist/eventmaker --metric "temp|celsius|float|0:72.1|3s"
```

To find out the endpoint to which your events are published see the [endpoints](#endpoints) section above


## cleanup 

To delete previously deployed fleet

```shell
bin/fleet-down
```

To delete hub and all of it's devices

> Note, this will delete the IoT Hub itself and all of its devices 

```shell
bin/hub-down
```


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](./LICENSE)


