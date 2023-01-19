import React from "react";
import { useRecoilValue } from "recoil";
import Title from "../../components/title/Title";
import { instance } from "../../recoil/atoms/instances";
import InstanceServicesTable from "./InstanceTable";

import "./Instances.scss";

const Instance = () => {
  const ins = useRecoilValue(instance);

  return (
    <main>
      <Title title={ins.name} subTitle={ins.description} />
      <React.Suspense fallback="Loading...">
        <InstanceServicesTable instance={ins} />
      </React.Suspense>
    </main>
  );
};

export default Instance;
