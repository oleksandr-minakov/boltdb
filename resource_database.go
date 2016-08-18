package boltdb

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"github.com/boltdb/bolt"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabase,
		Update: UpdateDatabase,
		Read:   ReadDatabase,
		Delete: DeleteDatabase,

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:		schema.TypeString,
				Required:	true,
				Default:	"BOLTDB_BUCKET",
			},
			"values": &schema.Schema{
				Type:		schema.TypeString,
				Optional:	true,
				Default:	"default_values",
			},
		},
	}
}

func CreateDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*bolt.DB)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(d.Get("bucket").(string)))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		log.Println("Bucket created: ", d.Get("bucket"))
		return nil
	})

	return nil
}

func UpdateDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*bolt.DB)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.Get("bucket").(string)))
		err := b.Put([]byte("KEY"), []byte(d.Get("values").(string)))
		log.Println("Bucket updated with values:", d.Get("values"))
		return err
	})

	return nil
}

func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*bolt.DB)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.Get("bucket").(string)))
		v := b.Get([]byte("KEY"))
		d.Set("values", string(v))
		log.Println("Read DB:", d.Get("bucket"))
		return nil
	})

	return nil
}

func DeleteDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*bolt.DB)

	db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(d.Get("bucket").(string)))
	})

	return nil
}
