import { SimpleForm } from "../../components/forms/Form";

type ProfileFormProps = {
  name: string;
  description: string;
};

const ProfileForm = ({
  data,
}: {
  create: boolean;
  data?: ProfileFormProps;
}) => {
  let fields = [
    {
      name: "name",
      help: "Name of the profile",
      fieldattrs: {
        type: "text",
        required: true,
        value: data?.name,
        placeholder: "Type a name for the profile",
      },
    },
    {
      name: "description",
      help: "Profile description",
      fieldattrs: {
        type: "text",
        required: false,
        value: data?.description,
        placeholder: "Add a short description for the profile",
        as: "textarea",
      },
    },
  ];

  return <SimpleForm page="Profiles" fields={fields} />;
};

export default ProfileForm;
