<template>
  <v-dialog v-model="dialog" width="80%" :theme="layoutTheme">
    <template v-slot:activator="{ props }">
      <v-btn :height="44" :width="44" v-bind="props" icon="mdi-help" color="primary" class="help" style="z-index: 9" />
    </template>

    <v-card>
      <v-toolbar title="Instructions">
        <template #append>
          <v-btn @click="dialog = false" icon="mdi-close"></v-btn>
        </template>
      </v-toolbar>
      <v-card-text>
        <h2 class="text-h6 mb-2 mt-3">Project Status</h2>

        <p>
          Please note that this project is currently in an alpha phase and is released as a preview build. If you encounter any problems, please feel free to
          file a bug report.
        </p>

        <h2 class="text-h6 mb-2 mt-3">This Playground supports the following kinds of rules</h2>

        <ul class="pl-4 my-2">
          <li>Validation Rules</li>
          <li>Mutation Rules</li>
          <li>Generation Rules</li>
          <li>Image Verification Rules</li>
        </ul>

        <p>It is not possible to access private registries (e.g., for image verification rules) or to use the clone capability in a generate rule.</p>

        <h2 class="text-h6 mb-2 mt-3">Context and Variables</h2>

        <p>
          It is currently not possible to add variables from external resources or do actual API calls. It is only possible to mock variables using the
          <code>variables</code>
          configuration in the context input.
        </p>

        <h2 class="text-h6 mb-2 mt-3">Multiple Manifests</h2>

        <p>It is supported to define multiple policies and/or resources as inputs. Context and variables will be shared for all executions.</p>

        <h2 class="text-h6 mb-2 mt-3">Load Manifests</h2>

        <ul class="pl-4 my-2 text-body-1">
          <li>The "File" Button loads a local YAML file as input</li>
          <li>
            The "URL" Button loads a manifest from an external URL
            <br />
            Example:
            <a href="https://raw.githubusercontent.com/kyverno/policies/main/best-practices/disallow-latest-tag/disallow-latest-tag.yaml" target="_blank">
              https://raw.githubusercontent.com/kyverno/policies/main/best-practices/disallow-latest-tag/disallow-latest-tag.yaml
            </a>
          </li>
        </ul>

        <h2 class="text-h6 mb-2 mt-3">Load / Save Functionality</h2>

        <ul class="pl-4 my-2 text-body-1">
          <li>
            Save new state: Your current state is saved in your local storage as a named profile, which means it is only available on your current machine and
            browser.
          </li>
          <li>Load a profile from your local storage, load the default profile to reset your current state to its defaults.</li>
          <li>Update a profile by selecting your current profile from the storage menu.</li>
        </ul>

        <h2 class="text-h6 mb-2 mt-3">Share</h2>

        <p>
          Share your current state with others by sending the created link. This link contains your state as a compressed base64 string, which means that this
          link will never expire or become invalid.
        </p>
      </v-card-text>
      <v-divider />
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { layoutTheme } from '@/config'
import { ref } from 'vue'

const dialog = ref<boolean>(false)
</script>

<style scoped>
.help {
  position: fixed;
  bottom: 45px;
  right: 185px;
}
</style>
