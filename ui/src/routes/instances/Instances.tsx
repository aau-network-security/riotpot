import "./Instances.scss";

import SearchBar from "../../components/searchbar/SearchBar";
import { InstancesTable } from "./InstancesTable";
import { InstancesUtils } from "./InstancesUtils";
import InstancesHeader from "./InstancesHeader";

const Page = "Instances";

/*
* This component is to create line charts.
* Enable and modify this component to add visual information regarding
* current instance traffic 

import { LineChart, Line } from "recharts";
function mockChart() {
  var n = [];
  while (n.length < 20) {
    var r = Math.floor(Math.random() * 100) + 1;
    n.push({ y: r });
  }

  return (
    <LineChart width={200} height={50} data={n}>
      <Line
        type="monotone"
        dataKey="y"
        stroke="#A4E18F"
        strokeWidth={3}
        dot={false}
      />
    </LineChart>
  );
}
*/

const Instances = () => {
  return (
    <main>
      <InstancesHeader view="Instances" />
      <InstancesUtils view={Page} />
      <SearchBar filter="instances" />
      <InstancesTable />
    </main>
  );
};

export default Instances;
