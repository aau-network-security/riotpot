import { CreateButton, UtilsBar } from "../../components/utils/Utils";
import { getPage } from "../../constants/globals";
import ServiceForm from "./ServiceForm";

export const ServicesUtils = ({ view }: { view: string }) => {
  // Utils buttons
  const page = getPage(view);
  const content = <ServiceForm create={true} />;
  const createbtn = (
    <CreateButton
      title="Service"
      icon={page?.icon}
      content={content}
    ></CreateButton>
  );
  const buttons = [createbtn];

  return <UtilsBar buttons={buttons} />;
};
