import { atom, Resetter, selectorFamily } from "recoil";
import { recoilPersist } from "recoil-persist";
import { Service } from "./services";
import uuid from "react-uuid";
import { ValidatorProps } from "./common";

const { persistAtom } = recoilPersist();

export type Profile = {
  id: string;
  name: string;
  description: string;
  services: Service[];
};

export const DefaultProfile = {
  id: "",
  name: "",
  description: "",
  services: [],
} as Profile;

export const profiles = atom({
  key: "profilesList",
  default: [] as Profile[],
  effects_UNSTABLE: [persistAtom],
});

export const profilesFilter = selectorFamily({
  key: "profile/default",
  get:
    (id: string | undefined) =>
    ({ get }) => {
      const profs = get(profiles);
      return profs.find((x: Profile) => x.id === id);
    },
  set:
    (id: string | undefined) =>
    ({ get, set }, newValue) => {
      const profs: any = get(profiles);
      const profInd = profs.findIndex((x: Profile) => x.id === id);

      let cp = [...profs];
      cp[profInd] = { ...cp[profInd], ...newValue };

      return set(profiles, cp);
    },
});

export const profileFormFields = atom({
  key: "profileFormFields",
  default: DefaultProfile as { [key: string]: any },
});

export const profileFormErrors = atom({
  key: "profileFormErrors",
  default: {} as { [key: string]: string },
});

export const profileFormField = selectorFamily({
  key: "profileFormField",
  get:
    (field?: string) =>
    ({ get }) =>
      !!field ? get(profileFormFields)[field] : get(profileFormFields),

  set:
    (field?: string) =>
    ({ set }, newValue) =>
      !!field
        ? // If a field is passed, set the value for the field
          set(profileFormFields, (prevState: any) => ({
            ...prevState,
            [field]: newValue,
          }))
        : // Otherwise, update the item values
          set(profileFormFields, (prev) => ({ ...prev, newValue })),
});

export const profileFormFieldErrors = selectorFamily({
  key: "profileFormFieldErrors",
  get:
    (field: string) =>
    ({ get }) =>
      get(profileFormErrors)[field],
});

export const AddProfileValidator = (
  profile: Profile,
  resetters: Resetter[],
  profileList: Profile[],
  setProfilesList: any,
  errors: any
) => {
  const err = "Invalid value";
  let errs = false;

  // Validate the fields
  profileList.forEach((curr) => {
    if (curr.name === profile.name) {
      errors((prev: any) => ({ ...prev, name: err }));
      errs = true;
    }
  });

  // Return if there are any errors
  if (errs) return;

  // Add the profile and reset the values
  let newprofile: Profile = {
    ...profile,
    id: uuid(),
  };
  setProfilesList((oldList: Profile[]) => [...oldList, newprofile]);

  // Reset the elements
  resetters.forEach((r) => r());
};

export const removeProfile = (
  id: string,
  profileList: any[],
  profileState: any
) => {
  profileState(profileList.filter((profile: Profile) => profile.id !== id));
};

export const profilesFilterBar = atom({
  key: "profilesFilterBar",
  default: "",
});

export const ProfileNameValidator = (props: ValidatorProps) => {
  props.query.forEach((q) => {
    if (q.name === props.element.name) {
      props.errors((prev: any) => ({ ...prev, name: "Invalid value" }));
    }
  });
};

/*
TODO: Unfinished. Take a look at this selector later on.

export const profilesFilter = selector({
  key: "profilesFilter",
  get: ({get}) => {
    const profilesList =  get(profiles);

    return profilesList.filter(profile => profile.id === id)
  },

  set: ({set}, newValue) => set(profiles, newValue),
})
*/

/**Validates the field value is unique in an array of objects */
/*
function ValidateIsFieldUniqueValue(field: any, value: any, query: Object[]) {
  for (let q in query) {
    if (q[field] === value) {
      return false;
    }
  }
  return true;
}
*/
