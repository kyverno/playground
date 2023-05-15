<template>
  <v-navigation-drawer
    :model-value="props.modelValue"
    temporary
    width="400"
    @update:modelValue="(event: boolean) => emit('update:modelValue', event)"
  >
    <template v-for="({ color, ...example }) in options.examples" :key="example.name">
      <v-list>
        <v-list-group :value="example.name">
          <template v-slot:activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-folder"
              :title="example.name"
              :class="color ? `text-${color}` : ''"
            ></v-list-item>
          </template>

          <template  v-if="example.subgroups">
            <template v-for="(subgroup, j) in example.subgroups" :key="j">
              <v-list-subheader>{{ subgroup.name }}</v-list-subheader>
              <v-list-item
                v-for="(policy, i) in subgroup.policies"
                :key="i"
                :title="formatTitle(policy)"
                @click="() => loadExample(policy.url || subgroup.url || example.url, policy)"
              />
            </template>
          </template>

          <template v-else>
            <v-list-item
              v-for="(policy, i) in example.policies"
              :key="i"
              :title="formatTitle(policy)"
              @click="() => loadExample(example.url, policy)"
            />
          </template>
        </v-list-group>
      </v-list>
    </template>

    <v-overlay :model-value="overlay" class="align-center justify-center">
      <v-progress-circular color="primary" indeterminate size="64"></v-progress-circular>
    </v-overlay>

    <template v-slot:append>
      <div class="pa-2">
        <v-btn flat color="primary" block @click="() => emit('update:modelValue', false)"
          >Close</v-btn
        >
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useConfig, Policy } from "../config";
import { ContextTemplate } from "@/assets/templates";

const { options } = useConfig()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
});

const emit = defineEmits(["update:modelValue", "select:example"]);
const overlay = ref<boolean>(false);

const formatTitle = (policy: Policy | string) => {
  if (typeof policy === 'object') {
    return policy.title
  }

  return policy.replaceAll('-', ' ')
}

const loadExample = async (url: string, policy: Policy | string) => {
  let folder: string = ''
  let contextPath: string | undefined

  if (typeof policy === 'object') {
    folder = policy.path
    contextPath = policy.contextPath
  } else {
    folder = policy
  }

  const name = folder.split('/').pop()

  try {
    const contextURL = contextPath ? `${contextPath}/${name}.yaml` : `${url}/${folder}/context.yaml`

    const promises = [
      fetch(`${url}/${folder}/${name}.yaml`).then((resp) => resp.text()),
      fetch(`${url}/${folder}/resource.yaml`).then((resp) => {
        if (resp.status === 404) {
          return fetch(`${url}/${folder}/resources.yaml`)
        }

        return resp
      }).then((resp) => resp.text()),
      fetch(contextURL).then((resp) => {
        if (resp.status === 404) {
          return ContextTemplate
        }

        return resp.text()
      }),
    ]

    overlay.value = true;
    const values = await Promise.all(promises);

    emit("select:example", values);
    emit("update:modelValue", false);
  } catch (err) {
    console.error(err);
  } finally {
    overlay.value = false;
  }
};
</script>
