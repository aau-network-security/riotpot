import { useState } from "react";
import { Button, Form, InputGroup, Row, Stack } from "react-bootstrap";
import { MutableSnapshot, Snapshot, useRecoilCallback } from "recoil";
import { instanceIds, instances } from "../../recoil/atoms/instances";
import { profiles } from "../../recoil/atoms/profiles";
import { services } from "../../recoil/atoms/services";
import stateData from "../../data/riotpot-state.json";

export function initState(snapshot: MutableSnapshot) {
  if (!stateData) return;

  snapshot.set(profiles, stateData.profiles);
  snapshot.set(services, stateData.services);

  snapshot.set(
    instanceIds,
    stateData.instances.map((inst: any) => inst.id)
  );

  stateData.instances.forEach((inst: any) =>
    snapshot.set(instances(inst.id), inst)
  );
}

async function processSnapshot(snapshot: Snapshot) {
  const persistedServices = await snapshot.getPromise(services);
  const persistedProfiles = await snapshot.getPromise(profiles);
  const persistedInstances = [];
  const persistedInstancesIds = await snapshot.getPromise(instanceIds);
  for (let instId of persistedInstancesIds) {
    const inst = await snapshot.getPromise(instances(instId));
    persistedInstances.push(inst);
  }

  const data = JSON.stringify({
    services: persistedServices,
    profiles: persistedProfiles,
    instances: persistedInstances,
  });

  localStorage.setItem("riotpot_storage", data);
  return data;
}

const LoadStateButton = () => {
  const [files, setFiles] = useState({});

  const loadState = useRecoilCallback(({ set }) => (data: any) => {
    set(services, data.services);
    set(profiles, data.profiles);

    const ids = data["instances"].map((inst: any) => inst.id);
    set(instanceIds, ids);

    for (let inst of data.instances) {
      set(instances(inst.id), inst);
    }
  });

  const handleChange = (event: any) => {
    const fileReader = new FileReader();
    const { files } = event.target;

    fileReader.readAsText(files[0], "UTF-8");
    fileReader.onloadend = (e) => {
      const content = e.target?.result;
      const toString = String(content);
      const toObj = JSON.parse(toString);
      setFiles(toObj);
    };
  };

  return (
    <Form.Group>
      <Form.Label>Upload</Form.Label>
      <InputGroup className="mb-3">
        <Form.Control type="file" onChange={handleChange} />
        <Button onClick={() => loadState(files)}>Upload</Button>
      </InputGroup>
      <Form.Text className="text-muted">
        Load the state from a JSON file
      </Form.Text>
    </Form.Group>
  );
};

const DownloadState = () => {
  const saveState = useRecoilCallback(({ snapshot }) => () => {
    processSnapshot(snapshot)
      .then((data) => {
        const blob = new Blob([data], { type: "text/plain" });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.download = "riotpot-state.json";
        link.href = url;
        link.click();
      })
      .catch((err) => {
        console.log(err);
      });
  });

  return (
    <Form.Group>
      <Form.Label>Download</Form.Label>
      <InputGroup className="mb-3">
        <Button onClick={saveState}>Download State</Button>
      </InputGroup>
      <Form.Text className="text-muted">
        Download the the current state to a JSON file
      </Form.Text>
    </Form.Group>
  );
};

export const SettingsState = () => {
  return (
    <Stack gap={5}>
      <Row>
        <h1>State</h1>
        <small>Manage the current state</small>
      </Row>
      <Row>
        <Form>
          <Stack gap={4}>
            <DownloadState />
            <LoadStateButton />
          </Stack>
        </Form>
      </Row>
    </Stack>
  );
};
