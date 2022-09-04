import React from "react";
import { useParams } from "react-router-dom";
import { useRecoilValue } from "recoil";
import Title from "../../components/title/Title";
import { instances } from "../../recoil/atoms/instances";
import InstanceServicesTable from "./InstanceTable";

import "./Instances.scss";

const Instance = () => {
  let { id } = useParams();

  const idN = id ? parseInt(id) : -1;
  const instance = useRecoilValue(instances(idN));

  return (
    <main>
      <Title title={instance.name} subTitle={instance.description} />
      <React.Suspense fallback="Loading...">
        <InstanceServicesTable instance={instance} />
      </React.Suspense>
    </main>
  );
};

export default Instance;
