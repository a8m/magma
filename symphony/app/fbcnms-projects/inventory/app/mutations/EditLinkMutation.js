/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  EditLinkMutationResponse,
  EditLinkMutationVariables,
} from './__generated__/EditLinkMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

export const mutation = graphql`
  mutation EditLinkMutation($input: EditLinkInput!) {
    editLink(input: $input) {
      id
      ...EquipmentPortsTable_link @relay(mask: false)
    }
  }
`;

export default (
  variables: EditLinkMutationVariables,
  callbacks?: MutationCallbacks<EditLinkMutationResponse>,
  updater?: (store: any) => void,
) => {
  const {onCompleted, onError} = callbacks || {};
  commitMutation(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
