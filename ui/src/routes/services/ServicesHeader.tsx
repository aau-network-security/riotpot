import Title from "../../components/title/Title";

const ServicesHeader = ({ view }: { view: string }) => {
  const subTitle: string = "List of registered services.";

  return <Title title={view} subTitle={subTitle} />;
};

export default ServicesHeader;
