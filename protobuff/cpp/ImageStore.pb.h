// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: ImageStore.proto

#ifndef PROTOBUF_INCLUDED_ImageStore_2eproto
#define PROTOBUF_INCLUDED_ImageStore_2eproto

#include <string>

#include <google/protobuf/stubs/common.h>

#if GOOGLE_PROTOBUF_VERSION < 3006001
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please update
#error your headers.
#endif
#if 3006001 < GOOGLE_PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_table_driven.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/inlined_string_field.h>
#include <google/protobuf/metadata.h>
#include <google/protobuf/message.h>
#include <google/protobuf/repeated_field.h>  // IWYU pragma: export
#include <google/protobuf/extension_set.h>  // IWYU pragma: export
#include <google/protobuf/unknown_field_set.h>
// @@protoc_insertion_point(includes)
#define PROTOBUF_INTERNAL_EXPORT_protobuf_ImageStore_2eproto 

namespace protobuf_ImageStore_2eproto {
// Internal implementation detail -- do not use these members.
struct TableStruct {
  static const ::google::protobuf::internal::ParseTableField entries[];
  static const ::google::protobuf::internal::AuxillaryParseTableField aux[];
  static const ::google::protobuf::internal::ParseTable schema[6];
  static const ::google::protobuf::internal::FieldMetadata field_metadata[];
  static const ::google::protobuf::internal::SerializationTable serialization_table[];
  static const ::google::protobuf::uint32 offsets[];
};
void AddDescriptors();
}  // namespace protobuf_ImageStore_2eproto
namespace ImageStore {
class ReadReq;
class ReadReqDefaultTypeInternal;
extern ReadReqDefaultTypeInternal _ReadReq_default_instance_;
class ReadResp;
class ReadRespDefaultTypeInternal;
extern ReadRespDefaultTypeInternal _ReadResp_default_instance_;
class RemoveReq;
class RemoveReqDefaultTypeInternal;
extern RemoveReqDefaultTypeInternal _RemoveReq_default_instance_;
class RemoveResp;
class RemoveRespDefaultTypeInternal;
extern RemoveRespDefaultTypeInternal _RemoveResp_default_instance_;
class StoreReq;
class StoreReqDefaultTypeInternal;
extern StoreReqDefaultTypeInternal _StoreReq_default_instance_;
class StoreResp;
class StoreRespDefaultTypeInternal;
extern StoreRespDefaultTypeInternal _StoreResp_default_instance_;
}  // namespace ImageStore
namespace google {
namespace protobuf {
template<> ::ImageStore::ReadReq* Arena::CreateMaybeMessage<::ImageStore::ReadReq>(Arena*);
template<> ::ImageStore::ReadResp* Arena::CreateMaybeMessage<::ImageStore::ReadResp>(Arena*);
template<> ::ImageStore::RemoveReq* Arena::CreateMaybeMessage<::ImageStore::RemoveReq>(Arena*);
template<> ::ImageStore::RemoveResp* Arena::CreateMaybeMessage<::ImageStore::RemoveResp>(Arena*);
template<> ::ImageStore::StoreReq* Arena::CreateMaybeMessage<::ImageStore::StoreReq>(Arena*);
template<> ::ImageStore::StoreResp* Arena::CreateMaybeMessage<::ImageStore::StoreResp>(Arena*);
}  // namespace protobuf
}  // namespace google
namespace ImageStore {

// ===================================================================

class ReadReq : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:ImageStore.ReadReq) */ {
 public:
  ReadReq();
  virtual ~ReadReq();

  ReadReq(const ReadReq& from);

  inline ReadReq& operator=(const ReadReq& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  ReadReq(ReadReq&& from) noexcept
    : ReadReq() {
    *this = ::std::move(from);
  }

  inline ReadReq& operator=(ReadReq&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const ReadReq& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const ReadReq* internal_default_instance() {
    return reinterpret_cast<const ReadReq*>(
               &_ReadReq_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  void Swap(ReadReq* other);
  friend void swap(ReadReq& a, ReadReq& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline ReadReq* New() const final {
    return CreateMaybeMessage<ReadReq>(NULL);
  }

  ReadReq* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<ReadReq>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const ReadReq& from);
  void MergeFrom(const ReadReq& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(ReadReq* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // string readKeyname = 1;
  void clear_readkeyname();
  static const int kReadKeynameFieldNumber = 1;
  const ::std::string& readkeyname() const;
  void set_readkeyname(const ::std::string& value);
  #if LANG_CXX11
  void set_readkeyname(::std::string&& value);
  #endif
  void set_readkeyname(const char* value);
  void set_readkeyname(const char* value, size_t size);
  ::std::string* mutable_readkeyname();
  ::std::string* release_readkeyname();
  void set_allocated_readkeyname(::std::string* readkeyname);

  // @@protoc_insertion_point(class_scope:ImageStore.ReadReq)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::internal::ArenaStringPtr readkeyname_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_ImageStore_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class ReadResp : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:ImageStore.ReadResp) */ {
 public:
  ReadResp();
  virtual ~ReadResp();

  ReadResp(const ReadResp& from);

  inline ReadResp& operator=(const ReadResp& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  ReadResp(ReadResp&& from) noexcept
    : ReadResp() {
    *this = ::std::move(from);
  }

  inline ReadResp& operator=(ReadResp&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const ReadResp& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const ReadResp* internal_default_instance() {
    return reinterpret_cast<const ReadResp*>(
               &_ReadResp_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    1;

  void Swap(ReadResp* other);
  friend void swap(ReadResp& a, ReadResp& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline ReadResp* New() const final {
    return CreateMaybeMessage<ReadResp>(NULL);
  }

  ReadResp* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<ReadResp>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const ReadResp& from);
  void MergeFrom(const ReadResp& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(ReadResp* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // bytes chunk = 1;
  void clear_chunk();
  static const int kChunkFieldNumber = 1;
  const ::std::string& chunk() const;
  void set_chunk(const ::std::string& value);
  #if LANG_CXX11
  void set_chunk(::std::string&& value);
  #endif
  void set_chunk(const char* value);
  void set_chunk(const void* value, size_t size);
  ::std::string* mutable_chunk();
  ::std::string* release_chunk();
  void set_allocated_chunk(::std::string* chunk);

  // @@protoc_insertion_point(class_scope:ImageStore.ReadResp)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::internal::ArenaStringPtr chunk_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_ImageStore_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class StoreReq : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:ImageStore.StoreReq) */ {
 public:
  StoreReq();
  virtual ~StoreReq();

  StoreReq(const StoreReq& from);

  inline StoreReq& operator=(const StoreReq& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  StoreReq(StoreReq&& from) noexcept
    : StoreReq() {
    *this = ::std::move(from);
  }

  inline StoreReq& operator=(StoreReq&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const StoreReq& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const StoreReq* internal_default_instance() {
    return reinterpret_cast<const StoreReq*>(
               &_StoreReq_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    2;

  void Swap(StoreReq* other);
  friend void swap(StoreReq& a, StoreReq& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline StoreReq* New() const final {
    return CreateMaybeMessage<StoreReq>(NULL);
  }

  StoreReq* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<StoreReq>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const StoreReq& from);
  void MergeFrom(const StoreReq& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(StoreReq* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // bytes chunk = 1;
  void clear_chunk();
  static const int kChunkFieldNumber = 1;
  const ::std::string& chunk() const;
  void set_chunk(const ::std::string& value);
  #if LANG_CXX11
  void set_chunk(::std::string&& value);
  #endif
  void set_chunk(const char* value);
  void set_chunk(const void* value, size_t size);
  ::std::string* mutable_chunk();
  ::std::string* release_chunk();
  void set_allocated_chunk(::std::string* chunk);

  // string memoryType = 2;
  void clear_memorytype();
  static const int kMemoryTypeFieldNumber = 2;
  const ::std::string& memorytype() const;
  void set_memorytype(const ::std::string& value);
  #if LANG_CXX11
  void set_memorytype(::std::string&& value);
  #endif
  void set_memorytype(const char* value);
  void set_memorytype(const char* value, size_t size);
  ::std::string* mutable_memorytype();
  ::std::string* release_memorytype();
  void set_allocated_memorytype(::std::string* memorytype);

  // @@protoc_insertion_point(class_scope:ImageStore.StoreReq)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::internal::ArenaStringPtr chunk_;
  ::google::protobuf::internal::ArenaStringPtr memorytype_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_ImageStore_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class StoreResp : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:ImageStore.StoreResp) */ {
 public:
  StoreResp();
  virtual ~StoreResp();

  StoreResp(const StoreResp& from);

  inline StoreResp& operator=(const StoreResp& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  StoreResp(StoreResp&& from) noexcept
    : StoreResp() {
    *this = ::std::move(from);
  }

  inline StoreResp& operator=(StoreResp&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const StoreResp& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const StoreResp* internal_default_instance() {
    return reinterpret_cast<const StoreResp*>(
               &_StoreResp_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    3;

  void Swap(StoreResp* other);
  friend void swap(StoreResp& a, StoreResp& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline StoreResp* New() const final {
    return CreateMaybeMessage<StoreResp>(NULL);
  }

  StoreResp* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<StoreResp>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const StoreResp& from);
  void MergeFrom(const StoreResp& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(StoreResp* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // string storeKeyname = 1;
  void clear_storekeyname();
  static const int kStoreKeynameFieldNumber = 1;
  const ::std::string& storekeyname() const;
  void set_storekeyname(const ::std::string& value);
  #if LANG_CXX11
  void set_storekeyname(::std::string&& value);
  #endif
  void set_storekeyname(const char* value);
  void set_storekeyname(const char* value, size_t size);
  ::std::string* mutable_storekeyname();
  ::std::string* release_storekeyname();
  void set_allocated_storekeyname(::std::string* storekeyname);

  // @@protoc_insertion_point(class_scope:ImageStore.StoreResp)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::internal::ArenaStringPtr storekeyname_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_ImageStore_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class RemoveReq : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:ImageStore.RemoveReq) */ {
 public:
  RemoveReq();
  virtual ~RemoveReq();

  RemoveReq(const RemoveReq& from);

  inline RemoveReq& operator=(const RemoveReq& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  RemoveReq(RemoveReq&& from) noexcept
    : RemoveReq() {
    *this = ::std::move(from);
  }

  inline RemoveReq& operator=(RemoveReq&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const RemoveReq& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const RemoveReq* internal_default_instance() {
    return reinterpret_cast<const RemoveReq*>(
               &_RemoveReq_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    4;

  void Swap(RemoveReq* other);
  friend void swap(RemoveReq& a, RemoveReq& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline RemoveReq* New() const final {
    return CreateMaybeMessage<RemoveReq>(NULL);
  }

  RemoveReq* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<RemoveReq>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const RemoveReq& from);
  void MergeFrom(const RemoveReq& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(RemoveReq* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // string remKeyname = 1;
  void clear_remkeyname();
  static const int kRemKeynameFieldNumber = 1;
  const ::std::string& remkeyname() const;
  void set_remkeyname(const ::std::string& value);
  #if LANG_CXX11
  void set_remkeyname(::std::string&& value);
  #endif
  void set_remkeyname(const char* value);
  void set_remkeyname(const char* value, size_t size);
  ::std::string* mutable_remkeyname();
  ::std::string* release_remkeyname();
  void set_allocated_remkeyname(::std::string* remkeyname);

  // @@protoc_insertion_point(class_scope:ImageStore.RemoveReq)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::internal::ArenaStringPtr remkeyname_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_ImageStore_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class RemoveResp : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:ImageStore.RemoveResp) */ {
 public:
  RemoveResp();
  virtual ~RemoveResp();

  RemoveResp(const RemoveResp& from);

  inline RemoveResp& operator=(const RemoveResp& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  RemoveResp(RemoveResp&& from) noexcept
    : RemoveResp() {
    *this = ::std::move(from);
  }

  inline RemoveResp& operator=(RemoveResp&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const RemoveResp& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const RemoveResp* internal_default_instance() {
    return reinterpret_cast<const RemoveResp*>(
               &_RemoveResp_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    5;

  void Swap(RemoveResp* other);
  friend void swap(RemoveResp& a, RemoveResp& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline RemoveResp* New() const final {
    return CreateMaybeMessage<RemoveResp>(NULL);
  }

  RemoveResp* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<RemoveResp>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const RemoveResp& from);
  void MergeFrom(const RemoveResp& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(RemoveResp* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // @@protoc_insertion_point(class_scope:ImageStore.RemoveResp)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_ImageStore_2eproto::TableStruct;
};
// ===================================================================


// ===================================================================

#ifdef __GNUC__
  #pragma GCC diagnostic push
  #pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// ReadReq

// string readKeyname = 1;
inline void ReadReq::clear_readkeyname() {
  readkeyname_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& ReadReq::readkeyname() const {
  // @@protoc_insertion_point(field_get:ImageStore.ReadReq.readKeyname)
  return readkeyname_.GetNoArena();
}
inline void ReadReq::set_readkeyname(const ::std::string& value) {
  
  readkeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:ImageStore.ReadReq.readKeyname)
}
#if LANG_CXX11
inline void ReadReq::set_readkeyname(::std::string&& value) {
  
  readkeyname_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:ImageStore.ReadReq.readKeyname)
}
#endif
inline void ReadReq::set_readkeyname(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  readkeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:ImageStore.ReadReq.readKeyname)
}
inline void ReadReq::set_readkeyname(const char* value, size_t size) {
  
  readkeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:ImageStore.ReadReq.readKeyname)
}
inline ::std::string* ReadReq::mutable_readkeyname() {
  
  // @@protoc_insertion_point(field_mutable:ImageStore.ReadReq.readKeyname)
  return readkeyname_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* ReadReq::release_readkeyname() {
  // @@protoc_insertion_point(field_release:ImageStore.ReadReq.readKeyname)
  
  return readkeyname_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void ReadReq::set_allocated_readkeyname(::std::string* readkeyname) {
  if (readkeyname != NULL) {
    
  } else {
    
  }
  readkeyname_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), readkeyname);
  // @@protoc_insertion_point(field_set_allocated:ImageStore.ReadReq.readKeyname)
}

// -------------------------------------------------------------------

// ReadResp

// bytes chunk = 1;
inline void ReadResp::clear_chunk() {
  chunk_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& ReadResp::chunk() const {
  // @@protoc_insertion_point(field_get:ImageStore.ReadResp.chunk)
  return chunk_.GetNoArena();
}
inline void ReadResp::set_chunk(const ::std::string& value) {
  
  chunk_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:ImageStore.ReadResp.chunk)
}
#if LANG_CXX11
inline void ReadResp::set_chunk(::std::string&& value) {
  
  chunk_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:ImageStore.ReadResp.chunk)
}
#endif
inline void ReadResp::set_chunk(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  chunk_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:ImageStore.ReadResp.chunk)
}
inline void ReadResp::set_chunk(const void* value, size_t size) {
  
  chunk_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:ImageStore.ReadResp.chunk)
}
inline ::std::string* ReadResp::mutable_chunk() {
  
  // @@protoc_insertion_point(field_mutable:ImageStore.ReadResp.chunk)
  return chunk_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* ReadResp::release_chunk() {
  // @@protoc_insertion_point(field_release:ImageStore.ReadResp.chunk)
  
  return chunk_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void ReadResp::set_allocated_chunk(::std::string* chunk) {
  if (chunk != NULL) {
    
  } else {
    
  }
  chunk_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), chunk);
  // @@protoc_insertion_point(field_set_allocated:ImageStore.ReadResp.chunk)
}

// -------------------------------------------------------------------

// StoreReq

// bytes chunk = 1;
inline void StoreReq::clear_chunk() {
  chunk_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& StoreReq::chunk() const {
  // @@protoc_insertion_point(field_get:ImageStore.StoreReq.chunk)
  return chunk_.GetNoArena();
}
inline void StoreReq::set_chunk(const ::std::string& value) {
  
  chunk_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:ImageStore.StoreReq.chunk)
}
#if LANG_CXX11
inline void StoreReq::set_chunk(::std::string&& value) {
  
  chunk_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:ImageStore.StoreReq.chunk)
}
#endif
inline void StoreReq::set_chunk(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  chunk_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:ImageStore.StoreReq.chunk)
}
inline void StoreReq::set_chunk(const void* value, size_t size) {
  
  chunk_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:ImageStore.StoreReq.chunk)
}
inline ::std::string* StoreReq::mutable_chunk() {
  
  // @@protoc_insertion_point(field_mutable:ImageStore.StoreReq.chunk)
  return chunk_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* StoreReq::release_chunk() {
  // @@protoc_insertion_point(field_release:ImageStore.StoreReq.chunk)
  
  return chunk_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void StoreReq::set_allocated_chunk(::std::string* chunk) {
  if (chunk != NULL) {
    
  } else {
    
  }
  chunk_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), chunk);
  // @@protoc_insertion_point(field_set_allocated:ImageStore.StoreReq.chunk)
}

// string memoryType = 2;
inline void StoreReq::clear_memorytype() {
  memorytype_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& StoreReq::memorytype() const {
  // @@protoc_insertion_point(field_get:ImageStore.StoreReq.memoryType)
  return memorytype_.GetNoArena();
}
inline void StoreReq::set_memorytype(const ::std::string& value) {
  
  memorytype_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:ImageStore.StoreReq.memoryType)
}
#if LANG_CXX11
inline void StoreReq::set_memorytype(::std::string&& value) {
  
  memorytype_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:ImageStore.StoreReq.memoryType)
}
#endif
inline void StoreReq::set_memorytype(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  memorytype_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:ImageStore.StoreReq.memoryType)
}
inline void StoreReq::set_memorytype(const char* value, size_t size) {
  
  memorytype_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:ImageStore.StoreReq.memoryType)
}
inline ::std::string* StoreReq::mutable_memorytype() {
  
  // @@protoc_insertion_point(field_mutable:ImageStore.StoreReq.memoryType)
  return memorytype_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* StoreReq::release_memorytype() {
  // @@protoc_insertion_point(field_release:ImageStore.StoreReq.memoryType)
  
  return memorytype_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void StoreReq::set_allocated_memorytype(::std::string* memorytype) {
  if (memorytype != NULL) {
    
  } else {
    
  }
  memorytype_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), memorytype);
  // @@protoc_insertion_point(field_set_allocated:ImageStore.StoreReq.memoryType)
}

// -------------------------------------------------------------------

// StoreResp

// string storeKeyname = 1;
inline void StoreResp::clear_storekeyname() {
  storekeyname_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& StoreResp::storekeyname() const {
  // @@protoc_insertion_point(field_get:ImageStore.StoreResp.storeKeyname)
  return storekeyname_.GetNoArena();
}
inline void StoreResp::set_storekeyname(const ::std::string& value) {
  
  storekeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:ImageStore.StoreResp.storeKeyname)
}
#if LANG_CXX11
inline void StoreResp::set_storekeyname(::std::string&& value) {
  
  storekeyname_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:ImageStore.StoreResp.storeKeyname)
}
#endif
inline void StoreResp::set_storekeyname(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  storekeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:ImageStore.StoreResp.storeKeyname)
}
inline void StoreResp::set_storekeyname(const char* value, size_t size) {
  
  storekeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:ImageStore.StoreResp.storeKeyname)
}
inline ::std::string* StoreResp::mutable_storekeyname() {
  
  // @@protoc_insertion_point(field_mutable:ImageStore.StoreResp.storeKeyname)
  return storekeyname_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* StoreResp::release_storekeyname() {
  // @@protoc_insertion_point(field_release:ImageStore.StoreResp.storeKeyname)
  
  return storekeyname_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void StoreResp::set_allocated_storekeyname(::std::string* storekeyname) {
  if (storekeyname != NULL) {
    
  } else {
    
  }
  storekeyname_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), storekeyname);
  // @@protoc_insertion_point(field_set_allocated:ImageStore.StoreResp.storeKeyname)
}

// -------------------------------------------------------------------

// RemoveReq

// string remKeyname = 1;
inline void RemoveReq::clear_remkeyname() {
  remkeyname_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& RemoveReq::remkeyname() const {
  // @@protoc_insertion_point(field_get:ImageStore.RemoveReq.remKeyname)
  return remkeyname_.GetNoArena();
}
inline void RemoveReq::set_remkeyname(const ::std::string& value) {
  
  remkeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:ImageStore.RemoveReq.remKeyname)
}
#if LANG_CXX11
inline void RemoveReq::set_remkeyname(::std::string&& value) {
  
  remkeyname_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:ImageStore.RemoveReq.remKeyname)
}
#endif
inline void RemoveReq::set_remkeyname(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  remkeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:ImageStore.RemoveReq.remKeyname)
}
inline void RemoveReq::set_remkeyname(const char* value, size_t size) {
  
  remkeyname_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:ImageStore.RemoveReq.remKeyname)
}
inline ::std::string* RemoveReq::mutable_remkeyname() {
  
  // @@protoc_insertion_point(field_mutable:ImageStore.RemoveReq.remKeyname)
  return remkeyname_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* RemoveReq::release_remkeyname() {
  // @@protoc_insertion_point(field_release:ImageStore.RemoveReq.remKeyname)
  
  return remkeyname_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void RemoveReq::set_allocated_remkeyname(::std::string* remkeyname) {
  if (remkeyname != NULL) {
    
  } else {
    
  }
  remkeyname_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), remkeyname);
  // @@protoc_insertion_point(field_set_allocated:ImageStore.RemoveReq.remKeyname)
}

// -------------------------------------------------------------------

// RemoveResp

#ifdef __GNUC__
  #pragma GCC diagnostic pop
#endif  // __GNUC__
// -------------------------------------------------------------------

// -------------------------------------------------------------------

// -------------------------------------------------------------------

// -------------------------------------------------------------------

// -------------------------------------------------------------------


// @@protoc_insertion_point(namespace_scope)

}  // namespace ImageStore

// @@protoc_insertion_point(global_scope)

#endif  // PROTOBUF_INCLUDED_ImageStore_2eproto
