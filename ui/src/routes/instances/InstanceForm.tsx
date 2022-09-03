import { SimpleForm } from "../../components/forms/Form";

type InstanceFormProps = {
  selector: any;
  errors: any;
  create: boolean;
  onSubmit: any;
};

export const InstanceFormFields = [
  {
    name: "name",
    help: "Name of the instance",
    fieldattrs: {
      type: "text",
      required: true,
      placeholder: "Type a name for the instance",
    },
  },
  {
    name: "description",
    help: "Instance description",
    fieldattrs: {
      type: "text",
      required: false,
      placeholder: "Add a short description for the instance",
      as: "textarea",
    },
  },
  {
    name: "host",
    help: "Host in where the instance can be reached. E.g.: localhost, 0.0.0.0, etc.",
    readOnly: true,
    fieldattrs: {
      type: "text",
      required: true,
      placeholder: "localhost",
    },
  },
];

const InstanceForm = (props: InstanceFormProps) => {
  // TODO: Perhaps repalce this in the future with an API call to get the fields and their structure
  let fields = [
    {
      name: "name",
      help: "Name of the instance",
      fieldattrs: {
        type: "text",
        required: true,
        placeholder: "Type a name for the instance",
      },
    },
    {
      name: "description",
      help: "Instance description",
      fieldattrs: {
        type: "text",
        required: false,
        placeholder: "Add a short description for the instance",
        as: "textarea",
      },
    },
    {
      name: "host",
      help: "Host in where the instance can be reached. E.g.: localhost, 0.0.0.0, etc.",
      fieldattrs: {
        type: "text",
        required: true,
        placeholder: "localhost",
      },
    },
  ];

  return (
    <SimpleForm
      onSubmit={props.onSubmit}
      page="Instances"
      fields={fields}
    ></SimpleForm>
  );
};

export default InstanceForm;
