// @generated by protoc-gen-es v1.3.0 with parameter "target=ts"
// @generated from file mgmt/v1alpha1/auth.proto (package mgmt.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64 } from "@bufbuild/protobuf";

/**
 * @generated from message mgmt.v1alpha1.GetNewAccessTokenRequest
 */
export class GetNewAccessTokenRequest extends Message<GetNewAccessTokenRequest> {
  /**
   * @generated from field: string client_id = 1;
   */
  clientId = "";

  /**
   * @generated from field: string refresh_token = 2;
   */
  refreshToken = "";

  constructor(data?: PartialMessage<GetNewAccessTokenRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "mgmt.v1alpha1.GetNewAccessTokenRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "client_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "refresh_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetNewAccessTokenRequest {
    return new GetNewAccessTokenRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetNewAccessTokenRequest {
    return new GetNewAccessTokenRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetNewAccessTokenRequest {
    return new GetNewAccessTokenRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetNewAccessTokenRequest | PlainMessage<GetNewAccessTokenRequest> | undefined, b: GetNewAccessTokenRequest | PlainMessage<GetNewAccessTokenRequest> | undefined): boolean {
    return proto3.util.equals(GetNewAccessTokenRequest, a, b);
  }
}

/**
 * @generated from message mgmt.v1alpha1.GetNewAccessTokenResponse
 */
export class GetNewAccessTokenResponse extends Message<GetNewAccessTokenResponse> {
  /**
   * @generated from field: string access_token = 1;
   */
  accessToken = "";

  /**
   * @generated from field: string refresh_token = 2;
   */
  refreshToken = "";

  /**
   * @generated from field: int64 expires_in = 3;
   */
  expiresIn = protoInt64.zero;

  /**
   * @generated from field: string scope = 4;
   */
  scope = "";

  /**
   * @generated from field: string id_token = 5;
   */
  idToken = "";

  /**
   * @generated from field: string token_type = 6;
   */
  tokenType = "";

  constructor(data?: PartialMessage<GetNewAccessTokenResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "mgmt.v1alpha1.GetNewAccessTokenResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "access_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "refresh_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "expires_in", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 4, name: "scope", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "id_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "token_type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetNewAccessTokenResponse {
    return new GetNewAccessTokenResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetNewAccessTokenResponse {
    return new GetNewAccessTokenResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetNewAccessTokenResponse {
    return new GetNewAccessTokenResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetNewAccessTokenResponse | PlainMessage<GetNewAccessTokenResponse> | undefined, b: GetNewAccessTokenResponse | PlainMessage<GetNewAccessTokenResponse> | undefined): boolean {
    return proto3.util.equals(GetNewAccessTokenResponse, a, b);
  }
}

/**
 * @generated from message mgmt.v1alpha1.GetAccessTokenRequest
 */
export class GetAccessTokenRequest extends Message<GetAccessTokenRequest> {
  /**
   * @generated from field: string client_id = 1;
   */
  clientId = "";

  /**
   * @generated from field: string code = 2;
   */
  code = "";

  /**
   * @generated from field: string redirect_uri = 3;
   */
  redirectUri = "";

  constructor(data?: PartialMessage<GetAccessTokenRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "mgmt.v1alpha1.GetAccessTokenRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "client_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "code", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "redirect_uri", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetAccessTokenRequest {
    return new GetAccessTokenRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetAccessTokenRequest {
    return new GetAccessTokenRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetAccessTokenRequest {
    return new GetAccessTokenRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetAccessTokenRequest | PlainMessage<GetAccessTokenRequest> | undefined, b: GetAccessTokenRequest | PlainMessage<GetAccessTokenRequest> | undefined): boolean {
    return proto3.util.equals(GetAccessTokenRequest, a, b);
  }
}

/**
 * @generated from message mgmt.v1alpha1.GetAccessTokenResponse
 */
export class GetAccessTokenResponse extends Message<GetAccessTokenResponse> {
  /**
   * @generated from field: string access_token = 1;
   */
  accessToken = "";

  /**
   * @generated from field: string refresh_token = 2;
   */
  refreshToken = "";

  /**
   * @generated from field: int64 expires_in = 3;
   */
  expiresIn = protoInt64.zero;

  /**
   * @generated from field: string scope = 4;
   */
  scope = "";

  /**
   * @generated from field: string id_token = 5;
   */
  idToken = "";

  /**
   * @generated from field: string token_type = 6;
   */
  tokenType = "";

  constructor(data?: PartialMessage<GetAccessTokenResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "mgmt.v1alpha1.GetAccessTokenResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "access_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "refresh_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "expires_in", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 4, name: "scope", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "id_token", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "token_type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetAccessTokenResponse {
    return new GetAccessTokenResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetAccessTokenResponse {
    return new GetAccessTokenResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetAccessTokenResponse {
    return new GetAccessTokenResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetAccessTokenResponse | PlainMessage<GetAccessTokenResponse> | undefined, b: GetAccessTokenResponse | PlainMessage<GetAccessTokenResponse> | undefined): boolean {
    return proto3.util.equals(GetAccessTokenResponse, a, b);
  }
}
