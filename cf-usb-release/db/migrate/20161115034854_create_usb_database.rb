class CreateUsbDatabase < ActiveRecord::Migration
  def self.up
  create_table "Config", id: false, force: true do |t|
    t.string "KEY",                    null: false
    t.string "VALUE",     limit: 1000
    t.string "COMPONENT"
  end

  create_table "Dials", primary_key: "Guid", force: true do |t|
    t.binary "Configuration"
    t.string "Plans_Guid",     limit: 36, null: false
    t.string "Instances_Guid", limit: 32, null: false
  end

  add_index "Dials", ["Instances_Guid"], name: "fk_Dials_Instances1_idx", using: :btree
  add_index "Dials", ["Plans_Guid"], name: "fk_Dials_Plans1_idx", using: :btree

  create_table "Instances", primary_key: "Guid", force: true do |t|
    t.string  "Name",      limit: 45
    t.string  "TargetURL", limit: 45
    t.string  "AuthKey",   limit: 45
    t.binary  "CaCert"
    t.boolean "SkipSSL"
  end

  create_table "Plans", primary_key: "Guid", force: true do |t|
    t.text    "Name",        limit: 255
    t.text    "Description", limit: 255
    t.boolean "Free"
    t.binary  "Metadata"
  end

  create_table "Services", primary_key: "Guid", force: true do |t|
    t.boolean "Bindable"
    t.binary  "DashboardClient"
    t.text    "Description",     limit: 255
    t.binary  "Metadata"
    t.text    "Name",            limit: 255
    t.boolean "PlanUpdateable"
    t.binary  "Tags"
    t.string  "Instances_Guid",  limit: 32,  null: false
  end

  add_index "Services", ["Instances_Guid"], name: "fk_Services_Instances1_idx", using: :btree

  end
end
