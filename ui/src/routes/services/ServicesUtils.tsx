import { useRecoilValue } from "recoil";
import { SimpleForm } from "../../components/forms/Form";
import { CreateButton, UtilsBar } from "../../components/utils/Utils";
import { getPage } from "../../constants/globals";
import { AddWithValidators, CreateItem } from "../../recoil/atoms/common";
import {
  DefaultService,
  serviceFormErrors,
  serviceFormFieldErrors,
  ServiceNameValidator,
  services,
  servicesFormFields,
} from "../../recoil/atoms/services";
import { ServiceFormFields } from "./ServiceForm";

const CreateService = () => {
  // Utils buttons
  const page = getPage("Services");

  var props: any = CreateItem({
    items: services,
    filter: servicesFormFields,
    errors: serviceFormErrors,
    validators: [ServiceNameValidator],
  });

  const onSubmit = (newElement: any) => {
    const updatedProps = {
      ...props,
      newElement: {
        ...DefaultService,
        ...newElement,
      },
    };

    AddWithValidators(updatedProps);
  };

  const defaultValues = useRecoilValue(servicesFormFields);

  const content = (
    <SimpleForm
      create={true}
      defaultValues={defaultValues}
      errors={serviceFormFieldErrors}
      onSubmit={onSubmit}
      page="Services"
      fields={ServiceFormFields}
    />
  );

  return (
    <CreateButton
      title="Service"
      icon={page?.icon}
      content={content}
    ></CreateButton>
  );
};

export const ServicesUtils = () => {
  const buttons = [<CreateService />];
  return <UtilsBar buttons={buttons} />;
};
