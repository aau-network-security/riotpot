import Title from "../../components/title/Title";

const InstancesHeader = ({ view }: { view: string }) => {
  const subTitle: string = "List of instances running.";
  return <Title title={view} subTitle={subTitle} />;
};

export default InstancesHeader;
