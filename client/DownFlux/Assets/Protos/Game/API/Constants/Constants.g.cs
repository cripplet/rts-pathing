// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: api/constants.proto
// </auto-generated>
#pragma warning disable 1591, 0612, 3021
#region Designer generated code

using pb = global::Google.Protobuf;
using pbc = global::Google.Protobuf.Collections;
using pbr = global::Google.Protobuf.Reflection;
using scg = global::System.Collections.Generic;
namespace DF.Game.API.Constants {

  /// <summary>Holder for reflection information generated from api/constants.proto</summary>
  public static partial class ConstantsReflection {

    #region Descriptor
    /// <summary>File descriptor for api/constants.proto</summary>
    public static pbr::FileDescriptor Descriptor {
      get { return descriptor; }
    }
    private static pbr::FileDescriptor descriptor;

    static ConstantsReflection() {
      byte[] descriptorData = global::System.Convert.FromBase64String(
          string.Concat(
            "ChNhcGkvY29uc3RhbnRzLnByb3RvEhJnYW1lLmFwaS5jb25zdGFudHMqTwoI",
            "TW92ZVR5cGUSFQoRTU9WRV9UWVBFX1VOS05PV04QABIVChFNT1ZFX1RZUEVf",
            "Rk9SV0FSRBABEhUKEU1PVkVfVFlQRV9SRVRSRUFUEAIqiQEKDkVudGl0eVBy",
            "b3BlcnR5EhsKF0VOVElUWV9QUk9QRVJUWV9VTktOT1dOEAASHAoYRU5USVRZ",
            "X1BST1BFUlRZX1BPU0lUSU9OEAESIAocRU5USVRZX1BST1BFUlRZX0FUVEFD",
            "S19USU1FUhACEhoKFkVOVElUWV9QUk9QRVJUWV9IRUFMVEgQAyqAAQoJQ3Vy",
            "dmVUeXBlEhYKEkNVUlZFX1RZUEVfVU5LTk9XThAAEhoKFkNVUlZFX1RZUEVf",
            "TElORUFSX01PVkUQARITCg9DVVJWRV9UWVBFX1NURVAQAhIUChBDVVJWRV9U",
            "WVBFX0RFTFRBEAMSFAoQQ1VSVkVfVFlQRV9USU1FUhAEKlgKCkVudGl0eVR5",
            "cGUSFwoTRU5USVRZX1RZUEVfVU5LTk9XThAAEhQKEEVOVElUWV9UWVBFX1RB",
            "TksQARIbChdFTlRJVFlfVFlQRV9FTlRJVFlfTElTVBACQixaEmdhbWUuYXBp",
            "LmNvbnN0YW50c6oCFURGLkdhbWUuQVBJLkNvbnN0YW50c2IGcHJvdG8z"));
      descriptor = pbr::FileDescriptor.FromGeneratedCode(descriptorData,
          new pbr::FileDescriptor[] { },
          new pbr::GeneratedClrTypeInfo(new[] {typeof(global::DF.Game.API.Constants.MoveType), typeof(global::DF.Game.API.Constants.EntityProperty), typeof(global::DF.Game.API.Constants.CurveType), typeof(global::DF.Game.API.Constants.EntityType), }, null, null));
    }
    #endregion

  }
  #region Enums
  /// <summary>
  /// MoveType represents a subtype of movement specified in a Move request.
  /// </summary>
  public enum MoveType {
    [pbr::OriginalName("MOVE_TYPE_UNKNOWN")] Unknown = 0,
    [pbr::OriginalName("MOVE_TYPE_FORWARD")] Forward = 1,
    [pbr::OriginalName("MOVE_TYPE_RETREAT")] Retreat = 2,
  }

  /// <summary>
  /// EntityProperty indicates the metric / property a curve represents.
  /// </summary>
  public enum EntityProperty {
    [pbr::OriginalName("ENTITY_PROPERTY_UNKNOWN")] Unknown = 0,
    [pbr::OriginalName("ENTITY_PROPERTY_POSITION")] Position = 1,
    [pbr::OriginalName("ENTITY_PROPERTY_ATTACK_TIMER")] AttackTimer = 2,
    [pbr::OriginalName("ENTITY_PROPERTY_HEALTH")] Health = 3,
  }

  /// <summary>
  /// CurveType indicates the interpolation method that should be used for the
  /// specified curve.
  /// </summary>
  public enum CurveType {
    [pbr::OriginalName("CURVE_TYPE_UNKNOWN")] Unknown = 0,
    /// <summary>
    /// TODO(minkezhang): Rename to LINEAR_POSITION
    /// </summary>
    [pbr::OriginalName("CURVE_TYPE_LINEAR_MOVE")] LinearMove = 1,
    [pbr::OriginalName("CURVE_TYPE_STEP")] Step = 2,
    [pbr::OriginalName("CURVE_TYPE_DELTA")] Delta = 3,
    [pbr::OriginalName("CURVE_TYPE_TIMER")] Timer = 4,
  }

  /// <summary>
  /// EntityType indicates the type of an object.
  /// </summary>
  public enum EntityType {
    [pbr::OriginalName("ENTITY_TYPE_UNKNOWN")] Unknown = 0,
    [pbr::OriginalName("ENTITY_TYPE_TANK")] Tank = 1,
    /// <summary>
    /// Server-only entity types.
    /// </summary>
    [pbr::OriginalName("ENTITY_TYPE_ENTITY_LIST")] EntityList = 2,
  }

  #endregion

}

#endregion Designer generated code