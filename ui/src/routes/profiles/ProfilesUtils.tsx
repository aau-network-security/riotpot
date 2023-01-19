import { CreateButton, UtilsBar } from "../../components/utils/Utils";
import { getPage } from "../../constants/globals";
import {
  DefaultProfile,
  profileFormErrors,
  profileFormFieldErrors,
  profileFormFields,
  ProfileNameValidator,
  profiles,
} from "../../recoil/atoms/profiles";
import { AddWithValidators, CreateItem } from "../../recoil/atoms/common";
import { SimpleForm } from "../../components/forms/Form";
import { ProfileFormFields } from "./ProfileForm";
import { useRecoilValue } from "recoil";

const CreateProfile = () => {
  const props: any = CreateItem({
    items: profiles,
    filter: profileFormFields,
    errors: profileFormErrors,
    validators: [ProfileNameValidator],
  });

  const onSubmit = (newElement: any) => {
    const updatedProps = {
      ...props,
      newElement: {
        ...DefaultProfile,
        ...newElement,
      },
    };

    AddWithValidators(updatedProps);
  };

  const defaultValues = useRecoilValue(profileFormFields);

  // Content of the form
  const content = (
    <SimpleForm
      create={true}
      defaultValues={defaultValues}
      errors={profileFormFieldErrors}
      onSubmit={onSubmit}
      page="Profiles"
      fields={ProfileFormFields}
    />
  );

  const page = getPage("Profiles");

  return <CreateButton title="Profile" icon={page?.icon} content={content} />;
};

export const ProfilesUtils = () => {
  // Utils buttons
  const buttons = [<CreateProfile key={0} />];
  return <UtilsBar buttons={buttons} />;
};
