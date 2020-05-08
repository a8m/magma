#!/usr/bin/env python3
# @generated AUTOGENERATED file. Do not Change!

from dataclasses import dataclass
from datetime import datetime
from gql.gql.datetime_utils import DATETIME_FIELD
from gql.gql.graphql_client import GraphqlClient
from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException
from functools import partial
from numbers import Number
from typing import Any, Callable, List, Mapping, Optional
from time import perf_counter
from dataclasses_json import DataClassJsonMixin

from ..fragment.property_type import PropertyTypeFragment, QUERY as PropertyTypeFragmentQuery
from ..input.add_location_type import AddLocationTypeInput


QUERY: List[str] = PropertyTypeFragmentQuery + ["""
mutation AddLocationTypeMutation($input: AddLocationTypeInput!) {
  addLocationType(input: $input) {
    id
    name
    propertyTypes {
      ...PropertyTypeFragment
    }
  }
}

"""]

@dataclass
class AddLocationTypeMutation(DataClassJsonMixin):
    @dataclass
    class AddLocationTypeMutationData(DataClassJsonMixin):
        @dataclass
        class LocationType(DataClassJsonMixin):
            @dataclass
            class PropertyType(PropertyTypeFragment):
                pass

            id: str
            name: str
            propertyTypes: List[PropertyType]

        addLocationType: LocationType

    data: AddLocationTypeMutationData

    @classmethod
    # fmt: off
    def execute(cls, client: GraphqlClient, input: AddLocationTypeInput) -> AddLocationTypeMutationData.LocationType:
        # fmt: off
        variables = {"input": input}
        try:
            start_time = perf_counter()
            response_text = client.call(''.join(set(QUERY)), variables=variables)
            res = cls.from_json(response_text).data
            elapsed_time = perf_counter() - start_time
            client.reporter.log_successful_operation("AddLocationTypeMutation", variables, elapsed_time)
            return res.addLocationType
        except OperationException as e:
            raise FailedOperationException(
                client.reporter,
                e.err_msg,
                e.err_id,
                "AddLocationTypeMutation",
                variables,
            )