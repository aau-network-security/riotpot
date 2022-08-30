export const ProfileFormFields = [
  {
    name: "name",
    help: "Name of the profile",
    fieldattrs: {
      type: "text",
      required: true,
      placeholder: "Type a name for the profile",
    },
  },
  {
    name: "description",
    help: "Profile description",
    fieldattrs: {
      type: "text",
      required: false,
      placeholder: "Add a short description for the profile",
      as: "textarea",
    },
  },
];
