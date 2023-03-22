import React from "react";
import { useRecoilValue } from "recoil";
import Title from "../../components/title/Title";
import { instance } from "../../recoil/atoms/instances";
import InstanceServicesTable from "./InstanceTable";

import "./Instances.scss";
import { Row } from "react-bootstrap";

const Instance = () => {
  const ins = useRecoilValue(instance);

  return (
    <main>
      <Title title={ins.name} subTitle={ins.description} />
      {(() => {
        if (ins.profile) {
          return (
            <Row>
              <h5>{ins.profile.name}</h5>
              <small>{ins.profile.description}</small>
            </Row>
          );
        }
      })()}
      <React.Suspense fallback="Loading...">
        <InstanceServicesTable />
      </React.Suspense>
    </main>
  );
};

export default Instance;
