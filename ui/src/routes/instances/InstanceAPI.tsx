import { InteractionOptions, NetworkOptions } from "../../constants/globals";
import { Service } from "../../recoil/atoms/services";

export const fetchProxy = async (host: string) => {
  let response = [];
  try {
    response = await fetch("http://" + host + "/api/proxies/")
      .then((response) => response.json())
      // Map the content of the response
      .then((data) =>
        data.map((element: any) => {
          if ("service" in element) {
            let service = element["service"];

            // Parse the network and the interaction
            if ("network" in service)
              service["network"] = NetworkOptions.find(
                (x) => x.value === service["network"]
              );
            if ("interaction" in service)
              service["interaction"] = InteractionOptions.find(
                (x) => x.value === service["interaction"]
              );

            // Re-assign the service
            element["service"] = service;
          }

          return element;
        })
      );
  } catch (error) {
    console.log(error);
  }

  return response;
};

export const patchService = async (host: string, service: Service) => {
  try {
    const response = await fetch(
      "http://" + host + "/api/services/" + service.id + "/",
      {
        method: "POST",
        body: JSON.stringify({
          name: service.name,
          port: service.port,
          host: service.host,
        }),
        headers: {
          "Content-type": "application/json; charset=UTF-8",
        },
      }
    )
      .then((response) => response.json())
      .then((data) => {
        // Parse the network and the interaction
        if ("network" in data)
          data["network"] = NetworkOptions.find(
            (x) => x.value === data["network"]
          );

        if ("interaction" in data)
          data["interaction"] = InteractionOptions.find(
            (x) => x.value === data["interaction"]
          );

        return {
          ...service,
          ...data,
        };
      });

    return response;
  } catch (err) {
    console.log(err);
  }

  return service;
};

export const changeProxyPort = async (
  host: string,
  proxyId: string,
  port: number
) => {
  try {
    await fetch("http://" + host + "/api/proxies/" + proxyId + "/port", {
      method: "POST",
      body: JSON.stringify({
        port: port,
      }),
      headers: {
        "Content-type": "application/json; charset=UTF-8",
      },
    })
      .then((response) => response.json())
      .catch((err) => {
        console.log(err.message);
      });
  } catch (err) {
    console.log(err);
  }
};

export const changeProxyStatus = async (
  host: string,
  id: string,
  status: string,
  changeStatus: (status: string) => any
) => {
  try {
    await fetch("http://" + host + "/api/proxies/" + id + "/status", {
      method: "POST",
      body: JSON.stringify({
        status: status,
      }),
      headers: {
        "Content-type": "application/json; charset=UTF-8",
      },
    })
      .then((response) => response.json())
      .then((data) => {
        if (["running", "stopped"].includes(data.status))
          changeStatus(data.status);
      })
      .catch((err) => {
        console.log(err.message);
      });
  } catch (err) {
    console.log(err);
  }
};
