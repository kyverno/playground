import { ConfigTemplate, ContextTemplate, PolicyTemplate, ResourceTemplate } from "@/assets/templates";
import { reactive } from "vue";
import { useState, Inputs } from "@/composables";

export const inputs = reactive({
    diffResources: false,
    policy: PolicyTemplate,
    oldResource: '',
    resource: ResourceTemplate,
    context: ContextTemplate,
    config: ConfigTemplate,
})

export const reset = () => {
    inputs.diffResources = false,
    inputs.policy = PolicyTemplate
    inputs.oldResource = ''
    inputs.resource = ResourceTemplate
    inputs.context = ContextTemplate
    inputs.config = ConfigTemplate
}

export const setDefaults = () => {
    init({
        policy: PolicyTemplate,
        diffResources: false,
        resource: ResourceTemplate,
        context: ContextTemplate,
        config: ConfigTemplate,
    })
}

export const init = (values: Inputs) => {
    const state = useState()

    if (typeof values.policy === 'string') {
        state.policy.value = values.policy;
        inputs.policy = values.policy;
    }

    if (typeof values.resource === 'string') {
        state.resource.value = values.resource;
        inputs.resource = values.resource;
    }

    if (values.diffResources !== undefined) {
        inputs.diffResources = values.diffResources
    }

    if (!values.oldResource && values.resource && inputs.diffResources) {
        state.oldResource.value = values.resource;
        inputs.oldResource = values.resource;
    }

    if (typeof values.oldResource === 'string') {
        inputs.diffResources = true;
        state.oldResource.value = values.oldResource;
        inputs.oldResource = values.oldResource;
    }

    if (typeof values.context === 'string') {
        state.context.value = values.context;
        inputs.context = values.context;
    }

    if (typeof values.config === 'string') {
        state.config.value = values.config;
        inputs.config = values.config;
    }

    state.name.value = values.name || "";
}

export const update = (values: Inputs) => {
    if (values.policy) {
        inputs.policy = values.policy;
    }

    if (values.oldResource) {
        inputs.oldResource = values.oldResource;
    }

    if (values.resource) {
        inputs.resource = values.resource;
    }

    if (values.context) {
        inputs.context = values.context;
    }

    if (values.config) {
        inputs.config = values.config;
    }
}

export const populate = () => {
    const state = useState()

    if (state.policy.value) {
        inputs.policy = state.policy.value;
    }
    if (state.oldResource.value) {
        inputs.diffResources = true;
        inputs.oldResource = state.oldResource.value;
    }
    if (state.resource.value) {
        inputs.resource = state.resource.value;
    }
    if (state.context.value) {
        inputs.context = state.context.value;
    }
    if (state.config.value) {
        inputs.config = state.config.value;
    }
}
